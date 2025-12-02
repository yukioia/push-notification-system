package http

import (
	"net/http"
	"push-server/src/internal/http/dto"
	"push-server/src/internal/model"
	"push-server/src/internal/service"
	"push-server/src/internal/sse"
	"sync"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	subService *service.SubscriptionService
	sseManager *sse.Manager
	clientsMu  sync.RWMutex
}

func NewHandler(subService *service.SubscriptionService, sseManager *sse.Manager) *Handler {
	return &Handler{
		subService: subService,
		sseManager: sseManager,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/subscribe", h.Subscribe)
		api.POST("/unsubscribe", h.Unsubscribe)
		api.GET("/subscriptions/:client_id", h.GetSubscriptions)
		api.POST("/publish", h.PublishMessage) // для теста рассылки
	}

	r.GET("/sse/:client_id", h.SSE)
}

// Subscribe — подписка на темы/теги
func (h *Handler) Subscribe(c *gin.Context) {
	var in dto.SubscriptionDTO
	if err := c.BindJSON(&in); err != nil || in.ClientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	sub := &model.Subscription{
		ClientID: in.ClientID,
		Topics:   in.Topics,
		Tags:     in.Tags,
	}
	if err := h.subService.Subscribe(sub); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "subscribed", "subscription": dto.FromDomain(sub)})
}

// Unsubscribe — отписка
func (h *Handler) Unsubscribe(c *gin.Context) {
	var in struct {
		ClientID string `json:"client_id"`
	}
	if err := c.BindJSON(&in); err != nil || in.ClientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_id required"})
		return
	}
	if err := h.subService.Unsubscribe(in.ClientID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.sseManager.Remove(in.ClientID)
	c.JSON(http.StatusOK, gin.H{"status": "unsubscribed"})
}

// GetSubscriptions — возвращает текущие подписки клиента
func (h *Handler) GetSubscriptions(c *gin.Context) {
	clientID := c.Param("client_id")
	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_id required"})
		return
	}
	subs, err := h.subService.GetByClientID(clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var out []dto.SubscriptionDTO
	for _, s := range subs {
		out = append(out, *dto.FromDomain(s))
	}
	c.JSON(http.StatusOK, out)
}

// SSE — открывает SSE поток для клиента
func (h *Handler) SSE(c *gin.Context) {
	clientID := c.Param("client_id")
	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_id required"})
		return
	}

	client := h.sseManager.Add(clientID)
	sse.SSEHandler(c.Writer, c.Request, client, h.sseManager)
}

// PublishMessage — для теста: отправка сообщений на клиентов
func (h *Handler) PublishMessage(c *gin.Context) {
	var in struct {
		Topic string `json:"topic"`
		Body  string `json:"body"`
	}
	if err := c.BindJSON(&in); err != nil || in.Topic == "" || in.Body == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	clientIDs, err := h.subService.GetClientsByTopic(in.Topic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, id := range clientIDs {
		_ = h.sseManager.SendTo(id, in.Body)
	}

	c.JSON(http.StatusOK, gin.H{"status": "message sent", "clients": clientIDs})
}
