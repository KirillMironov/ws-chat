package domain

const (
	ChatMessage      = "chat-message"
	ConnectedClients = "connected-clients"
)

type Message struct {
	Event   string  `json:"event"`
	Payload payload `json:"payload"`
}

type payload struct {
	Username string      `json:"username,omitempty"`
	Text     interface{} `json:"text"`
}
