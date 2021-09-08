package delivery

import (
	"github.com/KirillMironov/ws-chat/domain"
	"github.com/KirillMironov/ws-chat/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Handler struct {
	MessengerService service.Messenger
}

func NewHandler(ms service.Messenger) *Handler {
	return &Handler{MessengerService: ms}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("static/index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/connectToRoom", h.connectToRoom)
	return r
}

func (h *Handler) connectToRoom(c *gin.Context) {
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

	h.MessengerService.ConnectClient(&client, roomId)
}
