package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/knackwurstking/picow-led/picow"
	"github.com/lmittmann/tint"
)

const (
	// ErrorGeneric - every error not categorized
	ErrorGeneric = 1
	// ErrorArgs - invalid args given (non optional args)
	ErrorArgs = 2
	// ErrorInternal - something went wrong, this is a dev problem :)
	ErrorInternal = 10
	// ErrorServerError - something went wrong on the server side
	ErrorServerError = 15
	// ErrorUnderConstruction - feature not ready yet
	ErrorUnderConstruction = 100
)

var serverCache = &ServerCache{}

func main() {
	defer serverCache.Close()

	flags := NewFlags()
	flags.Read()

	level := slog.LevelInfo
	if flags.Debug {
		level = slog.LevelDebug
	}
	slog.SetDefault(
		slog.New(
			tint.NewHandler(
				os.Stderr,
				&tint.Options{
					AddSource:  true,
					TimeFormat: time.DateTime,
					Level:      level,
				},
			),
		),
	)

	slog.Debug("", "flags", flags)

	subs, err := flags.GetSubCommandArgs()
	if err != nil {
		slog.Error("Pasrsing flags failed ", "err", err)
		os.Exit(ErrorArgs)
	}

	for _, sub := range subs {
		// parse args for sub
		switch SubCommand(sub[0]) {
		case SubCommand_Run:
			subFlags, err := flags.ReadSubCommand(sub[1:])
			if err != nil {
				slog.Error("Parse ARGS failed", "command", sub[0], "err", err)
				os.Exit(ErrorArgs)
			}

			handleSubCommand_Run(
				flags.Addr,
				subFlags,
				newRequestFromARGS(subFlags.Args),
			).Wait()

		default:
			slog.Error("Ooops, subcommand not found", "command", sub[0])
			os.Exit(ErrorArgs)
		}
	}
}

func newRequestFromARGS(args []string) (req *picow.Request) {
	if len(args) < 3 {
		slog.Error("Wrong ARGS: <group> <command-type> <command> [<args> ...]")
		os.Exit(ErrorArgs)
	}

	group := picow.Group("")
	for _, g := range picow.Groups {
		if g == picow.Group(args[0]) {
			group = g
			break
		}
	}
	if group == "" {
		slog.Error("Group not exists!", "group", group)
		os.Exit(ErrorArgs)
	}

	_type := picow.Type("")
	for _, t := range picow.Types {
		if t == picow.Type(args[1]) {
			_type = t
			break
		}
	}
	if _type == "" {
		slog.Error("Type not exists!", "type", _type)
		os.Exit(ErrorArgs)
	}

	req = &picow.Request{
		ID:      0,
		Group:   picow.Group(args[0]),
		Type:    picow.Type(args[1]),
		Command: args[2],
		Args:    make([]string, 0),
	}

	return req
}

func handleSubCommand_Run(
	addrList AddrList,
	flags *FlagsSubCommand,
	request *picow.Request,
) *sync.WaitGroup {
	wg := sync.WaitGroup{}

	request.ID = flags.ID
	for _, a := range addrList {
		slog.Debug(
			"",
			"request", request,
			"address", a,
		)

		wg.Add(1)
		func(a string, wg *sync.WaitGroup) {
			defer wg.Done()

			server, err := serverCache.Get(a)
			if err != nil {
				slog.Warn(
					"Server connection for failed",
					"server", a,
					"err", err,
				)
				return
			}

			err = handleRequest(server, request, flags.PrettyPrint)
			if err != nil {
				slog.Warn(
					"Handle request failed",
					"server", a,
					"err", err,
				)
			}
		}(a, &wg)
	}

	return &wg
}

func handleRequest(
	server *picow.Server,
	request *picow.Request,
	prettyResponse bool,
) error {
	err := server.Send(request)
	if err != nil {
		return fmt.Errorf("request failed: %s", err.Error())
	}

	if request.ID == int(picow.IDNoResponse) {
		return nil
	}

	resp, err := server.GetResponse()
	if err != nil {
		return fmt.Errorf(
			"get response from \"%s\" failed: %s",
			server.GetAddr(), err.Error(),
		)
	}

	if resp.Error != "" {
		if resp.ID != 0 {
			err = fmt.Errorf(
				"id %d: %s: %s",
				resp.ID, server.GetAddr(), resp.Error,
			)
		} else {
			err = fmt.Errorf("%s: %s", server.GetAddr(), resp.Error)
		}
		return err
	}

	if resp.Data != nil {
		var data []byte
		if prettyResponse {
			data, err = json.MarshalIndent(resp.Data, "", "    ")
		} else {
			data, err = json.Marshal(resp.Data)
		}
		if err != nil {
			slog.Error(
				"Invalid json data from server",
				"server", server.GetAddr(),
				"resp.data", resp.Data,
			)
			os.Exit(ErrorServerError)
		}

		slog.Debug("", "resp", resp)
		fmt.Printf("%s\n", string(data))
	}

	return nil
}
