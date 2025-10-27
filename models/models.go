package models

import "net"

type DeviceID int64
type PinsID int64

// Addr contains an IP address and port number separated by a colon.
type Addr string

func (a Addr) SplitHostPort() (host string, port string) {
	host, port, _ = net.SplitHostPort(string(a)) // Ignore error
	return host, port
}
