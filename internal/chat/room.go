package chat

type Room struct {
	ID      int
	Name    string
	Clients map[string]*Client
}
