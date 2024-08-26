package main

import (
	"flag"
	"fmt"
	"log/slog"

	"github.com/knackwurstking/picow-led/picow"
)

// Flags holds all flag values
type Flags struct {
	Args  []string // Args containing all commandline args besides these already parsed
	Addr  AddrList // Addr containing the picow server addresses
	Debug bool     // Debug enables debugging messages
}

func NewFlags() *Flags {
	return &Flags{
		Args: make([]string, 0),
	}
}

// Read flags from args
func (f *Flags) Read() {
	flag.Var(&f.Addr, "addr", "picow device address (ip[:port] or hostname[:port])")
	flag.BoolVar(&f.Debug, "debug", f.Debug, "enable debug messages")

	flag.Parse()
	f.Args = flag.Args()
}

func (f *Flags) ReadSubCommand(name string, args []string) (*FlagsSubCommand, error) {
	cmd := flag.NewFlagSet(name, flag.ExitOnError)
	flags := NewFlagsSubCommand(name)

	cmd.IntVar(&flags.ID, "id", flags.ID, "changes the default id in use")
	cmd.BoolVar(&flags.PrettyPrint, "pretty-print", flags.PrettyPrint, "pretty prints response data")

	err := cmd.Parse(args)

	flags.Args = cmd.Args()
	slog.Debug("", "flags.Args", flags.Args)

	if flags.ID == int(picow.IDMotionEvent) && err == nil {
		err = fmt.Errorf("id \"%d\" not allowed", picow.IDMotionEvent)
	}

	return flags, err
}

func (f *Flags) GetSubCommandArgs() ([][]string, error) {
	subsArgs := make([][]string, 0)

	for _, arg := range f.Args {
		switch arg {
		case "get", "set":
			subsArgs = append(subsArgs, []string{arg})
		default:
			if len(subsArgs) == 0 {
				return subsArgs, fmt.Errorf("no sub command found")
			}

			subsArgs[len(subsArgs)-1] = append(subsArgs[len(subsArgs)-1], arg)
		}
	}

	return subsArgs, nil
}
