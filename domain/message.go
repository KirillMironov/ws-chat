package domain

const (
	ChatMessage      = "chat-message"
	ConnectedClients = "connected-clients"
)

type text interface {
	string | []string
}

type Message[T text] struct {
	Event      string `json:"event"`
	Payload[T] `json:"payload"`
}

type Payload[T text] struct {
	Username string `json:"username,omitempty"`
	Text     T      `json:"text"`
}
