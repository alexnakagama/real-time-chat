package server

import (
	"net/http"
	"strconv"

	"real-time-chat/internal/chat"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ServeWs(hub *chat.Hub, w http.ResponseWriter, r *http.Request) {
	roomIDStr := r.URL.Query().Get("room_id")
	userID := r.URL.Query().Get("user_id")
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		http.Error(w, "Invalid room_id", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade to websocket", http.StatusInternalServerError)
		return
	}

	client := &chat.Client{
		ID:     userID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		RoomID: roomID,
	}

	hub.Register <- client

	go client.WritePump()
	go client.ReadPump(hub)
}
