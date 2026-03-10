package env

import (
	"os"
	"path/filepath"

	"github.com/knackwurstking/ui"
)

const (
	DefaultVerbose       bool   = true
	DefaultServerAddress string = ":50835"
)

var (
	DBPath           string = os.Getenv("DB_PATH")
	ServerPathPrefix string = os.Getenv("SERVER_PATH_PREFIX")
	ServerAddress    string = os.Getenv("SERVER_ADDRESS")
	Verbose          bool   = os.Getenv("VERBOSE") == "true"

	//log *ui.Logger
)

func init() {
	if ServerAddress == "" {
		ServerAddress = DefaultServerAddress
	}
	if os.Getenv("VERBOSE") == "" {
		Verbose = DefaultVerbose
	}
	//log = NewLogger("env")
}

func NewLogger(name string) *ui.Logger {
	if Verbose {
		return ui.NewLoggerWithVerbose(name)
	} else {
		return ui.NewLogger(name)
	}
}

func Route(path string) string {
	return filepath.Join(ServerPathPrefix + path)
}
