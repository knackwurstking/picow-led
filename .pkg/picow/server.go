package picow

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	Name        string `json:"name"`
	Addr        string `json:"addr"`
	IsOffline   bool   `json:"isOffline"`
	IsConnected bool   `json:"-"`

	conn net.Conn
}

func NewServer(addr string) *Server {
	return &Server{
		Addr:      addr,
		IsOffline: false,
	}
}

func (s *Server) GetHost() string {
	return strings.Split(s.Addr, ":")[0]
}

func (s *Server) GetPort() (int, error) {
	as := strings.Split(s.Addr, ":")
	if len(as) != 2 {
		return 0, fmt.Errorf(
			"something is wrong with the server address: \"%s\"",
			s.Addr,
		)
	}

	return strconv.Atoi(as[1])
}

func (s *Server) Connect() error {
	d := net.Dialer{Timeout: time.Second * 5}
	c, err := d.Dial("tcp", s.Addr)
	if err != nil {
		s.IsOffline = true
		return err
	}

	s.conn = c
	s.IsConnected = true
	s.IsOffline = false

	return nil
}

func (s *Server) Close() error {
	s.IsConnected = false

	if s.conn == nil {
		return nil
	}

	return s.conn.Close()
}

func (s *Server) GetResponse() (*Response, error) {
	// check connection to the picow device
	if !s.IsConnected {
		return nil,
			fmt.Errorf(
				"not connected to server, run connect method first",
			)
	}

	// read data from client
	data := make([]byte, 0)
	chunk := make([]byte, 1)
	for {
		// read byte for byte and check for error
		s.conn.SetReadDeadline(time.Now().Add(time.Second * 5))
		n, err := s.conn.Read(chunk)
		if err != nil {
			s.IsOffline = true
			return nil, err
		}

		s.IsOffline = false

		// break on empty data
		if n == 0 {
			break
		}

		// checking for endbyte
		if chunk[0] == EndByte {
			break
		}

		// append chunk to data
		data = append(data, chunk...)
	}

	// check data
	if len(data) == 0 {
		return nil, fmt.Errorf("no data")
	}

	resp := Response{}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *Server) Send(req *Request) error {
	// Check connection to picow device
	if !s.IsConnected {
		return fmt.Errorf(
			"not connected to server, run connect method first",
		)
	}

	// Convert request to data
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	// Write data to client
	s.conn.SetWriteDeadline(time.Now().Add(time.Second * 5))
	n, err := s.conn.Write(append(data, EndByte))
	if err != nil {
		s.IsOffline = true
		return err
	} else if n == 0 {
		s.IsOffline = true
		return fmt.Errorf("no data written to \"%s\"", s.Addr)
	}

	s.IsOffline = false

	return nil
}
