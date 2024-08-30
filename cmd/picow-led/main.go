package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/knackwurstking/picow-led-server/pkg/picow"
	"github.com/knackwurstking/picow-led/cmd/picow-led/cache"
	"github.com/knackwurstking/picow-led/cmd/picow-led/errorcodes"
	"github.com/knackwurstking/picow-led/cmd/picow-led/flags"
)

var (
	serverCache = &cache.ServerCache{}
	prefixError = color.New(color.Bold, color.FgRed).Sprint("ERROR")
	prefixDebug = color.New(
		color.Bold, color.BgWhite, color.FgBlack,
	).Sprint("DEBUG")
)

func main() {
	defer serverCache.Close()

	flags := flags.NewFlags(serverCache)
	flags.Read()

	subs, err := flags.GetSubCommandArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"%s Parse sub commands failed: %s\n",
			prefixError, err)
		os.Exit(errorcodes.GetSubCommandArgs)
	}

	for _, sub := range subs {
		subFlags, err := flags.ReadSubCommand(sub[0], sub[1:])
		if err != nil {
			fmt.Fprintf(os.Stderr,
				"%s Parse sub command \"%s\" line flags failed: %s\n",
				prefixError, sub[0], err)
			os.Exit(errorcodes.ReadSubCommand)
		}

		wg := sync.WaitGroup{}
		for _, server := range serverCache.Data {
			wg.Add(1)
			go func(server *picow.Server) {
				defer wg.Done()

				if flags.Debug {
					fmt.Fprintf(os.Stderr,
						"%s Run \"%s %s\" Address %s\n",
						prefixDebug,
						subFlags.Flag.Name(), strings.Join(subFlags.Args, " "),
						server.GetAddr(),
					)
				}

				err := subFlags.Run(flags)
				if err != nil {
					fmt.Fprintf(os.Stderr,
						"%s %s\n", prefixError, err)
					os.Exit(errorcodes.Run)
				}
			}(server)
		}
		wg.Wait()
	}
}
