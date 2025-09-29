// main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	"encoding/json" 
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", ":8080", "http service address")

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true 
	},
}

type Client struct {
	conn     *websocket.Conn
	username string
}

type Message struct {
	Type      string   `json:"type"`
	Username  string   `json:"username,omitempty"`
	Content   string   `json:"content,omitempty"`
	Timestamp string   `json:"timestamp,omitempty"`
	Users     []string `json:"users,omitempty"`
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			h.sendUserList()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.conn.Close()
			}
			h.mu.Unlock()
			h.sendUserList()
		case msg := <-h.broadcast:
			h.mu.Lock()
			msg.Timestamp = time.Now().Format("15:04:05")
			messageData, err := json.Marshal(msg)
			if err != nil {
				log.Printf("marshal error: %v", err)
				h.mu.Unlock()
				continue
			}
			for client := range h.clients {
				err := client.conn.WriteMessage(websocket.TextMessage, messageData)
				if err != nil {
					log.Printf("write error: %v", err)
					client.conn.Close()
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) sendUserList() {
	h.mu.Lock()
	users := make([]string, 0, len(h.clients))
	for client := range h.clients {
		users = append(users, client.username)
	}
	h.mu.Unlock()

	userListMsg := Message{Type: "userlist", Users: users}
	messageData, err := json.Marshal(userListMsg)
	if err != nil {
		log.Printf("marshal error: %v", err)
		return
	}

	h.mu.Lock()
	for client := range h.clients {
		err := client.conn.WriteMessage(websocket.TextMessage, messageData)
		if err != nil {
			log.Printf("write error: %v", err)
			client.conn.Close()
			delete(h.clients, client)
		}
	}
	h.mu.Unlock()
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Read username as the first message
	_, usernameBytes, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}
	username := string(usernameBytes)
	if username == "" {
		conn.Close()
		return
	}
	client := &Client{conn: conn, username: username}
	hub.register <- client

	go func() {
		defer func() {
			hub.unregister <- client
		}()
		for {
			_, content, err := conn.ReadMessage()
			if err != nil {
				log.Printf("read error: %v", err)
				break
			}
			hub.broadcast <- Message{Type: "message", Username: client.username, Content: string(content)}
		}
	}()
}

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	fmt.Printf("Server starting on http://localhost:8080/ \n")
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}