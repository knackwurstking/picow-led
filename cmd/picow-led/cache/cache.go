package cache

import (
	"fmt"
	"log/slog"

	"github.com/knackwurstking/picow-led/picow"
)

type ServerCache struct {
	Data []*picow.Server
}

func (sc *ServerCache) Add(addr string) {
	sc.Data = append(sc.Data, picow.NewServer(addr))
}

func (sc *ServerCache) Get(addr string) (*picow.Server, error) {
	for _, server := range sc.Data {
		if server.GetAddr() == addr {
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
			slog.Warn(fmt.Sprintf("Close \"%s\" failed", server.GetAddr()), "err", err)
		}
	}
}
