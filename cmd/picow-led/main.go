package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/knackwurstking/picow-led/env"
	"github.com/knackwurstking/picow-led/services"
	"github.com/labstack/echo/v4"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	parseFlags()

	initializeLogging()
	r := initializeDatabase()

	startServer(r)
}

func parseFlags() {
	var logFormat string = string(env.Args.LogFormat)

	// Server Address
	flag.StringVar(&env.Args.Addr, "addr", env.Args.Addr, "Server address")

	// Server Path Prefix
	flag.StringVar(&env.Args.ServerPathPrefix, "path-prefix", env.Args.ServerPathPrefix, "Server path prefix")

	// Debug Flag
	flag.BoolVar(&env.Args.Debug, "debug", env.Args.Debug, "Enable debug mode")

	// Log Format: "text", "json"
	flag.StringVar(&logFormat, "log-format", logFormat, "Log format: text, json")

	// Get the required database path
	flag.StringVar(&env.Args.DatabasePath, "database-path", env.Args.DatabasePath, "Database path")

	// Custom Usage Message
	flag.Usage = func() {
		// TODO: List all exit codes and their meanings
		flag.PrintDefaults()
	}

	flag.Parse()

	switch logFormat {
	case "text":
		env.Args.LogFormat = env.LogFormatText
	case "json":
		env.Args.LogFormat = env.LogFormatJSON
	default:
		slog.Error("Invalid log format", "format", logFormat)
		os.Exit(env.ExitCodeInvalidLogFormat)
	}

	verifyDatabasePath()
}

func verifyDatabasePath() {
	// Verify the database path
	if env.Args.DatabasePath == "" {
		slog.Error("Database path is required")
		os.Exit(env.ExitCodeInvalidDatabasePath)
	}
}

func initializeLogging() {
	var level slog.Leveler
	if env.Args.Debug {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	var handler slog.Handler
	if env.Args.LogFormat == "text" {
		handler = slog.NewTextHandler(
			os.Stderr, &slog.HandlerOptions{Level: level},
		)
	} else {
		handler = slog.NewJSONHandler(
			os.Stderr, &slog.HandlerOptions{Level: level},
		)
	}

	slog.SetDefault(slog.New(handler))
}

func initializeDatabase() *services.Registry {
	slog.Debug("Initializing database", "database-path", env.Args.DatabasePath)

	sqlPath := fmt.Sprintf("%s", env.Args.DatabasePath)
	db, err := sql.Open("sqlite3", sqlPath)
	if err != nil {
		slog.Error("Failed to open database connection", "error", err)
		os.Exit(env.ExitCodeDatabaseConnection)
	}

	// Configure connection pool to prevent resource exhaustion
	db.SetMaxOpenConns(1)    // SQLite works best with single writer
	db.SetMaxIdleConns(1)    // Keep one connection alive
	db.SetConnMaxLifetime(0) // No maximum lifetime

	// Ping the database to verify the connection
	err = db.Ping()
	if err != nil {
		slog.Error("Failed to ping database", "error", err)
		os.Exit(env.ExitCodeDatabasePing)
	}

	return services.NewRegistry(db)
}

func startServer(r *services.Registry) {
	e := echo.New()

	// Initialize routes
	router(e, r)

	slog.Debug("Server started", "addr", env.Args.Addr, "server-path-prefix", env.Args.ServerPathPrefix)
	if err := e.Start(env.Args.Addr); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(env.ExitCodeServerStart)
	}
}
