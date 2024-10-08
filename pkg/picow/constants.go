package picow

import "time"

var (
	EndByte           = []byte("\n")
	SocketReadTimeout = time.Duration(time.Millisecond * 500) // 0.5s
)
