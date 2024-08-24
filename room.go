package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	clients map[*client]bool
	join    chan *client
	leave   chan *client
	forward chan []byte
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.receive)
		case msg := <-r.forward:
			for client := range r.clients {
				client.receive <- msg
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//First: Upgrade connection to a websocket
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("Error:", err)
		return
	}

	//Create a new instance of a client
	client := &client{
		socket:  socket,
		receive: make(chan []byte, messageBufferSize),
		room:    r,
	}
	//add the client to the room list ( room  map)
	r.join <- client

	//if connection closes for some reason, client leaves the room.
	defer func() { r.leave <- client }()

	//This handler will never exit unless the user closes browser, or some error occurs.
	//go routine so write and read concurrently
	go client.write()
	client.read()
}
