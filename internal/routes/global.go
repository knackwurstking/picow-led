package routes

type Global struct {
	ServerPathPrefix string
	Version          string
}

func (g Global) Array(v ...any) []any {
	return v
}
