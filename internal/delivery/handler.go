package delivery

import (
	"github.com/KirillMironov/ws-chat/domain"
	"github.com/KirillMironov/ws-chat/internal/ports"
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

func (h *Handler) InitRoutes() {
	http.Handle("/", http.FileServer(http.Dir("../static")))
	http.HandleFunc("/connect", h.connect)
}

func (h *Handler) connect(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	roomId := r.URL.Query().Get("roomId")
	if len(username) == 0 || len(roomId) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Info("not enough query params")
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
