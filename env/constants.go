package env

const (
	DefaultServerAddress    = ":50835"
	DefaultServerPathPrefix = "/picow-led"
)

const (
	ExitCodeServerStart         = 2
	ExitCodeInvalidLogFormat    = 3
	ExitCodeDatabaseConnection  = 4
	ExitCodeDatabasePing        = 5
	ExitCodeInvalidDatabasePath = 6
	ExitCodeInvalidFlags        = 8
)

const (
	CommandServer Command = "server"
)

var (
	Args *ArgsData
)

type Command string

type LogFormat string

const (
	LogFormatText LogFormat = "text"
	LogFormatJSON LogFormat = "json"
)

type ArgsData struct {
	Addr             string
	ServerPathPrefix string
	Debug            bool
	LogFormat        LogFormat
	DatabasePath     string
	Command          Command
}

func init() {
	Args = &ArgsData{
		Addr:             DefaultServerAddress,
		ServerPathPrefix: DefaultServerPathPrefix,
		LogFormat:        LogFormatJSON,
	}
}
