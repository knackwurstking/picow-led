package control

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/knackwurstking/picow-led/models"
)

const (
	EndByte byte = byte('\n')
)

var (
	controlMutexes = map[models.Addr]*sync.Mutex{}
)

type PicoW struct {
	*models.Device

	Conn      net.Conn
	connected models.Addr
}

// NewPicoW creates a new instance of the PicoW struct.
func NewPicoW(device *models.Device) *PicoW {
	return &PicoW{
		Device: device,
	}
}

// Write sends data to the Picow device. It first connects to the device if not already connected, then appends a newline character to the data and writes it.
func (p *PicoW) Write(request []byte) (n int, err error) {
	if p.Conn == nil {
		return 0, fmt.Errorf("not connected")
	}

	slog.Debug("Write data to device",
		"device_id", p.ID, "device_name", p.Name, "device_addr", p.Addr,
		"data", string(request))

	n, err = p.Conn.Write(p.EndByte(request))
	if err != nil {
		return 0, err
	}

	return n, nil
}

// Read reads up to len(response) bytes from the Picow device. It returns an error if not connected or no data is read.
func (p *PicoW) Read(response []byte) (n int, err error) {
	if p.Conn == nil {
		return 0, fmt.Errorf("not connected")
	}

	n, err = p.Conn.Read(response)
	if err != nil {
		return n, err
	}

	return n, nil
}

// ReadAll reads data from the Picow device until a newline character is encountered. It returns an error if not connected or no data is read.
func (p *PicoW) ReadAll() (data []byte, err error) {
	if p.Conn == nil {
		return nil, fmt.Errorf("not connected")
	}

	slog.Debug("Read data from device",
		"device_id", p.ID, "device_name", p.Name, "device_addr", p.Addr)

	buffer := bytes.NewBuffer(make([]byte, 0))
	chunk := make([]byte, 1)
	for {
		p.Conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		n, err := p.Read(chunk)
		if err != nil {
			return nil, err
		}
		if n == 0 {
			return nil, fmt.Errorf("no data read")
		}

		if bytes.Contains(chunk, []byte{EndByte}) {
			break
		}

		buffer.Write(chunk)
	}

	return buffer.Bytes(), nil
}

// Connect establishes a connection to the Picow device. If already connected, it returns nil.
func (p *PicoW) Connect() error {
	if p.Conn != nil {
		return nil
	}

	p.connected = p.Addr
	controlMutexes[p.connected] = &sync.Mutex{}

	dialer := net.Dialer{
		Timeout: time.Duration(time.Second * 5),
	}

	conn, err := dialer.Dial("tcp", string(p.Addr))
	if err != nil {
		return err
	}
	p.Conn = conn
	controlMutexes[p.connected].Lock()

	return nil
}

// Close closes the connection to the Picow device.
func (p *PicoW) Close() error {
	if p.Conn != nil {
		if err := p.Conn.Close(); err != nil {
			return err
		}

		p.Conn = nil
		if m, ok := controlMutexes[p.connected]; ok {
			m.Unlock()
			delete(controlMutexes, p.connected)
		}
		p.connected = ""
	}

	return nil
}

// EndByte returns the data with the end byte appended, only if not already present. The default end byte is a newline character.
func (p *PicoW) EndByte(data []byte) []byte {
	if len(data) == 0 || data[len(data)-1] != EndByte {
		return append(data, EndByte)
	}
	return data
}

var _ io.Writer = (*PicoW)(nil)
var _ io.Reader = (*PicoW)(nil)

// Error types returned:
//
//   - ErrNotConnected: Indicates that the connection to the Picow device is not established.
//   - ErrNoData: Indicates that no data was read from the Picow device.
