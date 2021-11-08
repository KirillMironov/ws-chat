package domain

type Message struct {
	Event   string  `json:"event"`
	Payload payload `json:"payload"`
}

type payload struct {
	Username string      `json:"username,omitempty"`
	Text     interface{} `json:"text"`
}
