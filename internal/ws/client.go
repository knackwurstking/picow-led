package ws

import "golang.org/x/net/websocket"

type Client struct {
	Conn *websocket.Conn
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		Conn: conn,
	}
}

func (c *Client) Send(data []byte) error {
	return websocket.Message.Send(c.Conn, data)
}
