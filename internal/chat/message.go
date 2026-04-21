package chat

type Message struct {
	SenderID string `json:"sender_id"`
	RoomID   int    `json:"room_id"`
	Content  string `json:"content"`
	Time     int64  `json:"time"`
}
