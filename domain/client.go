package domain

import "github.com/gorilla/websocket"

type Client struct {
	Username string
	RoomId   string
	Conn     *websocket.Conn
}
