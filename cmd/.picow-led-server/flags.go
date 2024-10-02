package main

import (
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/lmittmann/tint"
)

type Flags struct {
	Debug      bool
	Version    bool
	Port       int
	JSONLogger bool

	fmtHeading *color.Color
}

func NewFlags(port int) *Flags {
	return &Flags{
		Port:       port,
		fmtHeading: color.New(color.Bold, color.Underline, color.FgWhite),
	}
}

func (f *Flags) Parse() {
	flag.BoolVar(
		&f.JSONLogger, "json-logger", f.JSONLogger,
		"Use a JSON logger",
	)

	flag.BoolVar(
		&f.Debug, "debug", f.Debug,
		"Enable debugging logs",
	)

	flag.BoolVar(
		&f.Version, "version", f.Version,
		"Print bckend and frontend version and exit",
	)

	flag.IntVar(&f.Port, "port", f.Port, "Change the default port")

	flag.Usage = func() {
		f.fmtHeading.Fprintf(os.Stderr, "Options\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	f.setupLogger()
}

func (f *Flags) setupLogger() {
	var handler slog.Handler

	var level slog.Leveler
	if flags.Debug {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	if flags.JSONLogger {
		handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			AddSource: true,
			Level:     level,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				for _, g := range groups {
					if g == "request" || g == "response" {
						if a.Key == "time" ||
							a.Key == "query" ||
							a.Key == "referer" ||
							a.Key == "host" {
							return slog.Attr{}
						}
					}
				}

				return a
			},
		})
	} else {
		handler = tint.NewHandler(os.Stderr, &tint.Options{
			AddSource:  true,
			Level:      level,
			TimeFormat: time.UnixDate,
			ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
				for _, g := range groups {
					if g == "request" || g == "response" {
						if attr.Key == "time" ||
							attr.Key == "query" ||
							attr.Key == "referer" ||
							attr.Key == "host" {
							return slog.Attr{}
						}
					}
				}

				return attr
			},
		})
	}

	slog.SetDefault(slog.New(handler))
}
