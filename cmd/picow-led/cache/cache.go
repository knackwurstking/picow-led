package cache

import (
	"fmt"
	"log/slog"

	"github.com/knackwurstking/picow-led-server/pkg/picow"
)

type ServerCache struct {
	Data []*picow.Server
}

func (sc *ServerCache) Add(addr string) {
	sc.Data = append(sc.Data, picow.NewServer(addr))
}

func (sc *ServerCache) Get(addr string) (*picow.Server, error) {
	for _, server := range sc.Data {
		if server.Addr == addr {
			if server.IsConnected() {
				return server, nil
			} else {
				err := server.Connect()
				return server, err
			}
		}
	}

	server := picow.NewServer(addr)
	err := server.Connect()
	sc.Data = append(sc.Data, server)
	return server, err
}

func (sc *ServerCache) Close() {
	var err error
	for _, server := range sc.Data {
		err = server.Close()
		if err != nil {
			slog.Warn(
				fmt.Sprintf(
					"Close \"%s\" failed\n",
					server.Addr,
				),
				"err", err,
			)
		}
	}
}
