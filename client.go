package main

import (
	"github.com/gorilla/websocket"
)

type client struct {
	//socket is the websocket for this client
	socket *websocket.Conn
	//receive is a channel to get messages. array of bytes
	receive chan []byte
	//room is where client is chatting
	room *room
}

func (c *client) read() {
	defer c.socket.Close()

	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()

	for msg := range c.receive {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
