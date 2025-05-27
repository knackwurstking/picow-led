package main

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"picow-led/internal/database"
	"picow-led/internal/micro"
	"picow-led/internal/routes"
	"sync"
	"time"

	"github.com/SuperPaintman/nice/cli"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lmittmann/tint"
	"gopkg.in/yaml.v3"
)

var (
	ServerPathPrefix      = os.Getenv("SERVER_PATH_PREFIX")
	ServerAddress         = os.Getenv("SERVER_ADDR")
	Version               = "v0.1.0"
	ApiConfigPath         = "api.yaml"
	ApiConfigFallbackPath = ""              // See init()
	DBPath                = "./database.db" // TODO: Use a default system path for this (not the config directory)
)

type ConfigDevice struct {
	Addr string  `yaml:"addr"`
	Name string  `yaml:"name"`
	Pins []uint8 `yaml:"pins"`
}

type Config struct {
	Devices []ConfigDevice `yaml:"devices"`
}

func (c *Config) GetDataBaseDevices() []*database.Device {
	devices := []*database.Device{}

	if c.Devices == nil {
		return devices
	}

	wg := &sync.WaitGroup{}
	for _, d := range c.Devices {
		wg.Add(1)
		go func() {
			defer wg.Done()

			if d.Pins != nil {
				if len(d.Pins) > 0 {
					micro.SetPins(d.Addr, d.Pins)

					// TODO: Get color and create device and append to devices
				}
			}
		}()
	}
	wg.Wait()

	return devices
}

func init() {
	if d, err := os.UserConfigDir(); err == nil {
		ApiConfigFallbackPath = filepath.Join(d, "picow-led", "api.yaml")
	}

	// TODO: Get dir to data storage (used for the database)
}

func main() {
	app := cli.App{
		Name:  "picow-led",
		Usage: cli.Usage("PicoW LED server and client."),
		Commands: []cli.Command{
			{
				Name:  "server",
				Usage: cli.Usage("Start the Server"),
				Action: cli.ActionFunc(func(cmd *cli.Command) cli.ActionRunner {
					addr := cli.String(
						cmd, "addr",
						cli.WithShort("a"),
						cli.Usage("Set server address (<host>:<port>)"),
					)

					dbPath := cli.String(
						cmd, "db",
						cli.WithShort("d"),
						cli.Usage("Change database"),
					)

					*addr = ServerAddress
					*dbPath = DBPath

					return cliAction_Server(addr, dbPath)
				}),
			},
		},
		CommandFlags: []cli.CommandFlag{
			cli.HelpCommandFlag(),
			cli.VersionCommandFlag(Version),
		},
	}

	app.HandleError(app.Run())
}

func cliAction_Server(addr *string, dbPath *string) cli.ActionRunner {
	return func(cmd *cli.Command) error {
		e := echo.New()

		// Logger setup
		logger := slog.New(tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.DateTime,
			AddSource:  true,
		}))
		slog.SetDefault(logger)

		// Middleware
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "[${time_rfc3339}] ${status} ${method} ${path} (${remote_ip}) ${latency_human}\n",
			Output: os.Stderr,
		}))

		// Create and init the database
		db := database.NewDB(*dbPath)
		defer db.Close()

		// Load configuration and pass data to the database
		if config, err := loadConfig(ApiConfigPath, ApiConfigFallbackPath); err != nil {
			slog.Warn("Configuration", "error", err)
		} else {
			devices := config.GetDataBaseDevices()
			if err = db.Devices.Set(devices...); err != nil {
				slog.Error("Set devices to database", "error", err)
			}
		}

		// Register routes
		routesOptions := &routes.Options{
			ServerPathPrefix: ServerPathPrefix,
			DB:               db,
		}

		routes.Register(e, routesOptions)

		return e.Start(*addr)
	}
}

func loadConfig(paths ...string) (*Config, error) {
	config := &Config{}

	for _, path := range paths {
		absPath, err := filepath.Abs(path)
		if err != nil {
			absPath = path
		}
		f, err := os.Open(absPath)
		if err != nil {
			continue
		}
		d, err := io.ReadAll(f)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(d, config)
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}
