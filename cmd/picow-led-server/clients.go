package main

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ReadTimeout  time.Time
	WriteTimeout time.Time
	Conn         *websocket.Conn
	Chan         chan any

	wg     *sync.WaitGroup
	cancel chan any
	done   chan any

	Pulse time.Duration
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		ReadTimeout:  time.Now().Add(time.Millisecond * 2500),
		WriteTimeout: time.Now().Add(time.Millisecond * 1000),
		Conn:         conn,
		Chan:         make(chan any),

		wg:     &sync.WaitGroup{},
		cancel: make(chan any),
		done:   make(chan any),

		Pulse: time.Millisecond * 2500,
	}
}

func (c *Client) StartHeartBeat() {
	go func() {
		for {
			select {
			case <-c.cancel:
				return
			default:
				time.Sleep(c.Pulse)

				c.Conn.SetWriteDeadline(c.WriteTimeout)
				err := c.Conn.WriteMessage(websocket.BinaryMessage, []byte("ping"))
				if err != nil {
					c.done <- nil
					return
				}

				c.Conn.SetReadDeadline(c.ReadTimeout)
				_, p, err := c.Conn.ReadMessage()
				if (err == nil && string(p) != "pong") || err != nil {
					c.done <- nil
					return
				}
			}
		}
	}()
}

func (c *Client) StopHeartBeat() {
	c.cancel <- nil
	c.wg.Wait()
}

func (c *Client) Done() chan any {
	return c.done
}
