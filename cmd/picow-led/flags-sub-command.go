package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/knackwurstking/picow-led/picow"
)

type FlagsSubCommand struct {
	Name        string
	Args        []string
	ID          int
	PrettyPrint bool
}

func NewFlagsSubCommand(name string) *FlagsSubCommand {
	return &FlagsSubCommand{
		Name: name,
	}
}

func (fsc *FlagsSubCommand) Run(flags *Flags) error {
	switch picow.Type(fsc.Name) {
	case picow.TypeGet:
		return fsc.get(flags)
	case picow.TypeSet:
		return fsc.set(flags)
	}

	return nil
}

func (fsc *FlagsSubCommand) get(flags *Flags) error {
	r := fsc.request(picow.TypeGet)
	wg := sync.WaitGroup{}
	r.ID = fsc.ID

	for _, a := range flags.Addr {
		slog.Debug("", "request", r, "address", a)
		wg.Add(1)
		fsc.send(a, r, &wg)
	}

	return nil
}

func (fsc *FlagsSubCommand) set(flags *Flags) error {
	r := fsc.request(picow.TypeSet)
	wg := sync.WaitGroup{}
	r.ID = fsc.ID

	for _, a := range flags.Addr {
		slog.Debug("", "request", r, "address", a)
		wg.Add(1)
		fsc.send(a, r, &wg)
	}

	return nil
}

func (fsc *FlagsSubCommand) request(t picow.Type) *picow.Request {
	if len(fsc.Args) < 3 {
		slog.Error("Missing ARGS: <group> <command> [<args> ...]")
		os.Exit(ErrorArgs)
	}

	group := picow.Group("")
	for _, g := range picow.Groups {
		if g == picow.Group(fsc.Args[0]) {
			group = g
			break
		}
	}

	if group == "" {
		slog.Error("Group not exists!", "group", group)
		os.Exit(ErrorArgs)
	}

	return &picow.Request{
		ID:      0,
		Group:   picow.Group(fsc.Args[0]),
		Type:    t,
		Command: fsc.Args[1],
		Args:    make([]string, 0),
	}
}

func (fsc *FlagsSubCommand) send(addr string, r *picow.Request, wg *sync.WaitGroup) error {
	defer wg.Done()

	server, err := serverCache.Get(addr)
	if err != nil {
		return err
	}

	err = server.Send(r)
	if err != nil {
		return err
	}

	if r.ID == int(picow.IDNoResponse) {
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
		if fsc.PrettyPrint {
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
