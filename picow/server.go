package picow

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// Server will handle all communication to a picow device
type Server struct {
	conn        net.Conn
	addr        string
	isConnected bool
}

// NewServer will create a new Server object
func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

// GetHost of the current picow device
func (s *Server) GetHost() string {
	return strings.Split(s.addr, ":")[0]
}

// GetPort of the current picow device, returns zero on an error
func (s *Server) GetPort() (int, error) {
	as := strings.Split(s.addr, ":")
	if len(as) != 2 {
		return 0, fmt.Errorf(
			"something is wrong with the server address: \"%s\"",
			s.addr,
		)
	}

	return strconv.Atoi(as[1])
}

// GetAddr returns the <host>:<port>
func (s *Server) GetAddr() string {
	return s.addr
}

// Checks if the connection to the server is still up
func (s *Server) IsConnected() bool {
	return s.isConnected
}

// Connect to picow device socket, uses "tcp"
func (s *Server) Connect() error {
	c, err := net.Dial("tcp", s.addr)
	if err != nil {
		return err
	}

	s.conn = c
	s.isConnected = true

	return nil
}

func (s *Server) Close() error {
	if s.conn == nil {
		return nil
	}

	return s.conn.Close()
}

// GetResponse from the picow device
func (s *Server) GetResponse() (*Response, error) {
	// check connection to the picow device
	if !s.isConnected {
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
		s.conn.SetReadDeadline(time.Now().Add(time.Second))
		n, err := s.conn.Read(chunk)
		if err != nil {
			return nil, err
		}

		// break on empty data
		if n == 0 {
			break
		}

		// checking for endbyte
		if chunk[0] == DefaultEndByte {
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

// Send a request to the picow
func (s *Server) Send(req *Request) error {
	// Check connection to picow device
	if !s.isConnected {
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
	s.conn.SetWriteDeadline(time.Now().Add(time.Second))
	n, err := s.conn.Write(append(data, DefaultEndByte))
	if err != nil {
		return err
	} else if n == 0 {
		return fmt.Errorf("no data written to \"%s\"", s.addr)
	}

	return nil
}
