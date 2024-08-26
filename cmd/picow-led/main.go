package main

import (
	"log/slog"
	"os"
	"time"

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
		subFlags, err := flags.ReadSubCommand(sub[0], sub[1:])
		if err != nil {
			slog.Error("Parse ARGS failed", "command", sub[0], "err", err)
			os.Exit(ErrorArgs)
		}

		subFlags.Run(flags)
	}
}
