package clients

import (
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	EmitType EmitType
	Conn     *websocket.Conn
	Chan     chan any

	done chan any
}

func NewClient(eventType EmitType, conn *websocket.Conn) *Client {
	return &Client{
		EmitType: eventType,
		Conn:     conn,
		Chan:     make(chan any, 2),
		done:     make(chan any, 2),
	}
}

func (c *Client) StartHeartBeat() (exit chan any) {
	exit = make(chan any, 2)
	go func() {
		for {
			select {
			case <-exit:
				return
			default:
				time.Sleep(time.Millisecond * 2500)

				c.Conn.SetWriteDeadline(time.Now().Add(time.Millisecond * 1000))
				if err := c.Conn.WriteMessage(websocket.BinaryMessage, []byte("ping")); err != nil {
					c.done <- nil
					return
				}

				c.Conn.SetReadDeadline(time.Now().Add(time.Millisecond * 2500))
				if _, p, err := c.Conn.ReadMessage(); (err == nil && string(p) != "pong") || err != nil {
					c.done <- nil
					return
				}
			}
		}
	}()

	return exit
}

func (c *Client) Done() chan any {
	return c.done
}
