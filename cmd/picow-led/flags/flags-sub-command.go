package flags

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"
	"sync"

	"github.com/knackwurstking/picow-led/cmd/picow-led/cache"
	"github.com/knackwurstking/picow-led/picow"
)

type FlagsSubCommand struct {
	serverCache *cache.ServerCache
	Flag        *flag.FlagSet
	Args        []string
	ID          int
	PrettyPrint bool
}

func NewFlagsSubCommand(flagSet *flag.FlagSet, sc *cache.ServerCache) *FlagsSubCommand {
	return &FlagsSubCommand{
		serverCache: sc,
		Flag:        flagSet,
	}
}

func (fsc *FlagsSubCommand) Run(flags *Flags) error {
	switch picow.Type(fsc.Flag.Name()) {
	case picow.TypeGet:
		return fsc.get(flags)
	case picow.TypeSet:
		return fsc.set(flags)
	}

	return nil
}

func (fsc *FlagsSubCommand) get(flags *Flags) error {
	r, err := fsc.request(picow.TypeGet)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	r.ID = fsc.ID

	for _, a := range flags.Addr {
		wg.Add(1)
		err = fsc.send(a, r, &wg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (fsc *FlagsSubCommand) set(flags *Flags) error {
	r, err := fsc.request(picow.TypeSet)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	r.ID = fsc.ID

	for _, a := range flags.Addr {
		wg.Add(1)
		err = fsc.send(a, r, &wg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (fsc *FlagsSubCommand) request(t picow.Type) (*picow.Request, error) {
	if len(fsc.Args) < 2 {
		return nil, fmt.Errorf("don't know about this command: \"%s\"", strings.Join(fsc.Args, " "))
	}

	group := picow.Group("")
	for _, g := range picow.Groups {
		if g == picow.Group(fsc.Args[0]) {
			group = g
			break
		}
	}

	if group == "" {
		return nil, fmt.Errorf("sub command group \"%s\" not exists", group)
	}

	return &picow.Request{
		ID:      0,
		Group:   picow.Group(fsc.Args[0]),
		Type:    t,
		Command: fsc.Args[1],
		Args:    make([]string, 0),
	}, nil
}

func (fsc *FlagsSubCommand) send(addr string, r *picow.Request, wg *sync.WaitGroup) error {
	defer wg.Done()

	server, err := fsc.serverCache.Get(addr)
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
			return fmt.Errorf("Invalid json data from server \"%s\": %+v", server.GetAddr(), resp.Data)
		}

		fmt.Printf("%s\n", string(data))
	}

	return nil
}
