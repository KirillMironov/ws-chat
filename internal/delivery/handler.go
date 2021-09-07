package delivery

import (
	"encoding/json"
	"github.com/KirillMironov/ws-chat/domain"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	rooms = make(map[string]domain.Room)
)

func InitRoutes() *gin.Engine {
	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("static/index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/connectToRoom", connectToRoom)
	return r
}

func connectToRoom(c *gin.Context) {
	username := c.Query("username")
	roomId := c.Query("roomId")
	if len(username) == 0 || len(roomId) == 0 {
		log.Println("not enough query params")
		c.Status(http.StatusBadRequest)
		return
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	client := domain.Client{
		Username: username,
		Conn:     ws,
	}

	if _, ok := rooms[roomId]; ok {
		rooms[roomId].Clients[client] = true
	} else {
		createNewRoom(&client, roomId)
		go messageWriter(roomId, rooms[roomId].Broadcast)
	}

	log.Printf("new client: '%s', roomId: '%s'", client.Username, roomId)
	go messageReader(&client, roomId, rooms[roomId].Broadcast)
}

func createNewRoom(client *domain.Client, roomId string) {
	msg := make(chan domain.Message)
	room := domain.Room{
		Id:        roomId,
		Clients:   make(map[domain.Client]bool),
		Broadcast: msg,
	}
	room.Clients[*client] = true
	rooms[roomId] = room
}

func messageReader(client *domain.Client, roomId string, messages chan<- domain.Message) {
	defer func() {
		_ = client.Conn.Close()
		delete(rooms[roomId].Clients, *client)
		if len(rooms[roomId].Clients) == 0 {
			delete(rooms, roomId)
			log.Printf("deleted empty room: '%s'", roomId)
		}
		log.Printf("closed connection with client: '%s', roomId: '%s'", client.Username, roomId)
	}()

	for {
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		messages <- domain.Message{
			Username: client.Username,
			Text:     string(p),
		}
	}
}

func messageWriter(roomId string, messages <-chan domain.Message) {
	for {
		select {
		case msg := <-messages:
			for client := range rooms[roomId].Clients {
				js, err := json.Marshal(msg)
				if err != nil {
					log.Println(err)
					return
				}
				log.Printf("sending message: '%s' to client: '%s'", js, client.Username)

				err = client.Conn.WriteMessage(1, js)
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}
