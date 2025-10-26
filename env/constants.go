package env

const (
	DefaultServerAddress    = ":50835"
	DefaultServerPathPrefix = "/picow-led"
)

const (
	ExitCodeServerStart      = 2
	ExitCodeInvalidLogFormat = 3
)

var (
	Args *ArgsData
)

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
}

func init() {
	Args = &ArgsData{
		Addr:             DefaultServerAddress,
		ServerPathPrefix: DefaultServerPathPrefix,
		LogFormat:        LogFormatJSON,
	}
}
