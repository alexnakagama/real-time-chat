package chat

import "encoding/json"

type Hub struct {
	Rooms      map[int]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[int]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.handleRegister(client)
		case client := <-h.Unregister:
			h.handleUnregister(client)
		case message := <-h.Broadcast:
			h.handleBroadcast(message)
		}
	}
}

func (h *Hub) handleRegister(client *Client) {
	room, ok := h.Rooms[client.RoomID]
	if !ok {
		room = &Room{
			ID:      client.RoomID,
			Clients: make(map[string]*Client),
		}
		h.Rooms[client.RoomID] = room
	}
	room.Clients[client.ID] = client
}

func (h *Hub) handleUnregister(client *Client) {
	if room, ok := h.Rooms[client.RoomID]; ok {
		if _, exists := room.Clients[client.ID]; exists {
			delete(room.Clients, client.ID)
			close(client.Send)
			if len(room.Clients) == 0 {
				delete(h.Rooms, client.RoomID)
			}
		}
	}
}

func (h *Hub) handleBroadcast(message *Message) {
	data, err := json.Marshal(message)
	if err != nil {
		return
	}
	if room, ok := h.Rooms[message.RoomID]; ok {
		for _, client := range room.Clients {
			select {
			case client.Send <- data:
			default:
				close(client.Send)
				delete(room.Clients, client.ID)
			}
		}
	}
}
