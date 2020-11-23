package ws

import "log"

type Hub struct {
	clients    map[*Client]bool
	private    chan []byte
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Println(client.Username, "上线")
			i := 0
			for k, v := range h.clients {
				log.Println(k.Username, "->", v)
				i++
			}
			log.Println("一共有", i, "个连接")
		case client := <-h.unregister:
			h.clients[client] = false
			delete(h.clients, client)
			log.Println(client.Username, "下线")
		case message := <-h.broadcast:

			for client := range h.clients {

				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
