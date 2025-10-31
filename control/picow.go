package control

import (
	"bytes"
	"io"
	"net"
	"time"

	"github.com/knackwurstking/picow-led/models"
)

type PicoW struct {
	*models.ResolvedDevice

	Conn net.Conn
}

func NewPicoW(device *models.ResolvedDevice) *PicoW {
	return &PicoW{
		ResolvedDevice: device,
	}
}

func (p *PicoW) Write(request []byte) (n int, err error) {
	p.Connect()

	n, err = p.Conn.Write(p.EndByte(request))
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (p *PicoW) Read(response []byte) (n int, err error) {
	if p.Conn == nil {
		return 0, ErrNotConnected
	}

	n, err = p.Conn.Read(response)
	if err != nil {
		return n, err
	}

	return n, nil
}

func (p *PicoW) ReadAll() (data []byte, err error) {
	if p.Conn == nil {
		return nil, ErrNotConnected
	}

	buffer := bytes.NewBuffer(make([]byte, 0))
	chunk := make([]byte, 1)
	for {
		p.Conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		n, err := p.Read(chunk)
		if err != nil {
			return nil, err
		}
		if n == 0 {
			return nil, ErrNoData
		}

		if bytes.Contains(chunk, []byte{'\n'}) {
			break
		}

		buffer.Write(chunk)
	}

	return buffer.Bytes(), nil
}

func (p *PicoW) Connect() error {
	if p.Conn != nil {
		return nil
	}

	dialer := net.Dialer{
		Timeout: time.Duration(time.Second * 5),
	}

	var err error
	p.Conn, err = dialer.Dial("tcp", string(p.Addr))
	if err != nil {
		return err
	}

	return nil
}

func (p *PicoW) Close() error {
	if err := p.Conn.Close(); err != nil {
		return err
	}

	p.Conn = nil
	return nil
}

// EndByte returns the data with the end byte appended, only if not already present, newline will be used as end byte here
func (p *PicoW) EndByte(data []byte) []byte {
	if len(data) == 0 || data[len(data)-1] != '\n' {
		return append(data, '\n')
	}
	return data
}

var _ io.Writer = &PicoW{}
var _ io.Reader = &PicoW{}
