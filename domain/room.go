package domain

type Room struct {
	Id        string
	Clients   map[Client]bool
	Broadcast chan Message
}
