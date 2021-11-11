package v1

import (
	"github.com/KirillMironov/ws-chat/domain"
	"github.com/KirillMironov/ws-chat/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

type Handler struct {
	clientsService ports.ClientsService
	logger         ports.Logger
}

func NewHandler(clientsService ports.ClientsService, logger ports.Logger) *Handler {
	return &Handler{clientsService: clientsService, logger: logger}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
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
		c.Status(http.StatusBadRequest)
		h.logger.Info("not enough query params")
		return
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		h.logger.Error(err)
		return
	}

	client := &domain.Client{
		Username: username,
		RoomId:   roomId,
		Conn:     ws,
	}

	h.clientsService.Connect(client)
}
