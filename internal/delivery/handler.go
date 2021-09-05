package delivery

import (
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

	rooms    = make(map[string]map[domain.Client]bool)
	channels = make(map[string]chan string)
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
		RoomId:   roomId,
		Conn:     ws,
	}

	if rooms[roomId] == nil {
		rooms[roomId] = make(map[domain.Client]bool)
		rooms[roomId][client] = true
	} else {
		rooms[roomId][client] = true
	}

	if channels[roomId] == nil {
		msg := make(chan string)
		channels[roomId] = msg
		go messageWriter(roomId, channels[roomId])
	}

	log.Printf("new client '%s', roomId: '%s'", client.Username, client.RoomId)

	go messageReader(&client, channels[roomId])
}

func messageReader(client *domain.Client, messages chan string) {
	defer func() {
		_ = client.Conn.Close()
		delete(rooms[client.RoomId], *client)
		log.Printf("closed connection with '%s', roomId: '%s'", client.Username, client.RoomId)
	}()

	for {
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(string(p))

		messages <- string(p)
	}
}

func messageWriter(roomId string, messages chan string) {
	for {
		select {
		case msg := <-messages:
			for client := range rooms[roomId] {
				err := client.Conn.WriteMessage(1, []byte(msg))
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}
