package delivery

import (
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
	}

	clients = make(map[string]*websocket.Conn)
)

func InitRoutes() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("static/index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/createRoom", createRoom)
	//r.GET("/connectToRoom", connectToRoom)
	return r
}

func createRoom(c *gin.Context) {
	roomId := c.Query("roomId")
	if len(roomId) == 0 {
		log.Println("empty room id")
		c.Status(http.StatusBadRequest)
		return
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	msg := make(chan bool)

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	clients[roomId] = ws

	go messageReader(msg, ws)
	go messageWriter(msg, ws)
}

//func connectToRoom(c *gin.Context) {
//	roomId := c.Query("roomId")
//	if len(roomId) == 0 {
//		log.Println("empty room id")
//		c.Status(http.StatusBadRequest)
//		return
//	}
//
//	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
//
//	ws := clients[roomId]
//	if ws == nil {
//		log.Println("room doesn't exist")
//		c.Status(http.StatusBadRequest)
//		return
//	}
//
//	go messageReader(ws)
//}

func messageReader(msg chan bool, conn *websocket.Conn) {
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(string(p))

		if string(p) == "What time is it?" {
			msg <- true
		}
	}
}

func messageWriter(msg chan bool, conn *websocket.Conn) {
	for _ = range msg {
		if err := conn.WriteMessage(1, []byte(time.Now().String())); err != nil {
			log.Println(err)
			return
		}
	}
}
