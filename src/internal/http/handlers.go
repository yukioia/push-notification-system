package http

import (
	"net/http"
	"push-server/src/internal/model"
	"push-server/src/internal/service"
	"push-server/src/internal/sse"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.NotificationService
	sse *sse.Manager
}

func NewHandler(svc *service.NotificationService, sse *sse.Manager) *Handler {
	return &Handler{svc: svc, sse: sse}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.POST("/api/subscribe", h.Subscribe)
	r.POST("/api/unsubscribe", h.Unsubscribe)
	r.POST("/api/publish", h.Publish)
	r.GET("/sse/:client_id", h.SSEConnect)
}

func (h *Handler) Subscribe(c *gin.Context) {
	var s model.Subscription
	if err := c.BindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if s.ClientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_id required"})
		return
	}
	if err := h.svc.Subscribe(&s); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "subscribed"})
}

func (h *Handler) Unsubscribe(c *gin.Context) {
	var in struct {
		ClientID string `json:"client_id"`
	}
	if err := c.BindJSON(&in); err != nil || in.ClientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_id required"})
		return
	}
	if err := h.svc.Unsubscribe(in.ClientID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "unsubscribed"})
}

func (h *Handler) Publish(c *gin.Context) {
	var m model.Message
	if err := c.BindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if m.Topic == "" || m.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "topic and title required"})
		return
	}
	if err := h.svc.Publish(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "published"})
}

func (h *Handler) SSEConnect(c *gin.Context) {
	clientID := c.Param("client_id")
	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_id required"})
		return
	}
	client := h.sse.Add(clientID)
	sse.SSEHandler(c.Writer, c.Request, client, h.sse)
}
