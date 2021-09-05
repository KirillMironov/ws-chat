package delivery

import (
	"github.com/KirillMironov/ws-chat/domain"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	rooms = make(map[domain.Client]bool)
)

func InitRoutes() *gin.Engine {
	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("static/index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/connectToRoom", connectToRoom)
	go messageWriter()
	return r
}

func connectToRoom(c *gin.Context) {
	username := c.Query("username")
	if len(username) == 0 {
		log.Println("username wasn't provided")
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
	rooms[client] = true

	log.Printf("new client '%s'", client.Username)

	go messageReader(&client)
}

func messageReader(client *domain.Client) {
	defer func() {
		_ = client.Conn.Close()
		delete(rooms, *client)
		log.Printf("closed connection with '%s'", client.Username)
	}()

	for {
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(string(p))
	}
}

func messageWriter() {
	for {
		for client := range rooms {
			log.Printf("sending message to '%s'", client.Username)

			err := client.Conn.WriteMessage(1, []byte(time.Now().String()))
			if err != nil {
				log.Println(err)
				return
			}
		}

		time.Sleep(5 * time.Second)
	}
}
