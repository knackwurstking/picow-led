package flags

import (
	"flag"
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"github.com/knackwurstking/picow-led/cmd/picow-led/cache"
	"github.com/knackwurstking/picow-led/picow"
)

// AddrList contains strings "<ip/hostname>:<port>" for the picow devices to connect to
type AddrList []string

// String returns a string with all addresses
func (a AddrList) String() string {
	return strings.Join(a, ",")
}

// Set adds a new server
func (a *AddrList) Set(value string) error {
	matched, _ := regexp.MatchString("^.+:[0-9]+$", value)
	if !matched {
		// no match means we have to add the default port here
		value = fmt.Sprintf("%s:%d", strings.TrimRight(value, ":"), picow.DefaultPort)
	}

	*a = append(*a, value)

	return nil
}

// Flags holds all flag values
type Flags struct {
	serverCache *cache.ServerCache
	Args        []string
	Addr        AddrList
	Debug       bool
}

func NewFlags(sc *cache.ServerCache) *Flags {
	return &Flags{
		Args:        make([]string, 0),
		serverCache: sc,
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
	flags := NewFlagsSubCommand(name, f.serverCache)

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
