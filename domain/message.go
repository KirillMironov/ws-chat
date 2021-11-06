package domain

type Message struct {
	Event   string  `json:"event"`
	Payload payload `json:"payload"`
}

type payload struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}
