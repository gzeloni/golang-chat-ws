package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type client struct {
	conn *websocket.Conn
}

func (c *client) read() {
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		broadcast(msg)
	}
}

var clients = make(map[*client]bool)
var broadcastChannel = make(chan []byte)

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	go handleMessages()

	fmt.Println("Servidor iniciado na porta :8080")
	err := http.ListenAndServe(":2154", nil)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	c := &client{conn: conn}
	clients[c] = true

	c.read()
}

func handleMessages() {
	for {
		msg := <-broadcastChannel
		for client := range clients {
			err := client.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				delete(clients, client)
				client.conn.Close()
			}
		}
	}
}

func broadcast(msg []byte) {
	broadcastChannel <- msg
}
