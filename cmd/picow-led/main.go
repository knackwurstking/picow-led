package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/knackwurstking/picow-led/cmd/picow-led/cache"
	"github.com/knackwurstking/picow-led/cmd/picow-led/errorcodes"
	"github.com/knackwurstking/picow-led/cmd/picow-led/flags"
	"github.com/lmittmann/tint"
)

var serverCache = &cache.ServerCache{}

func main() {
	defer serverCache.Close()

	flags := flags.NewFlags(serverCache)
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

	subs, err := flags.GetSubCommandArgs()
	if err != nil {
		slog.Error("Pasrsing flags failed ", "err", err)
		os.Exit(errorcodes.Args)
	}

	for _, sub := range subs {
		subFlags, err := flags.ReadSubCommand(sub[0], sub[1:])
		if err != nil {
			slog.Error("Parse ARGS failed", "command", sub[0], "err", err)
			os.Exit(errorcodes.Args)
		}

		err = subFlags.Run(flags)
		if err != nil {
			slog.Error(
				fmt.Sprintf(
					"Failed to run %s",
					strings.Join(flags.Args, " "),
				),
			)
		}
	}
}
