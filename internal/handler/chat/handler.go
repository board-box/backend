package chat

import (
	"net/http"

	chatSvc "github.com/board-box/backend/internal/service/chat"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *chatSvc.Service
	authMW  func(c *gin.Context)
}

func New(service *chatSvc.Service, authMW func(c *gin.Context)) *Handler {
	return &Handler{service, authMW}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/chat")
	g.Use(h.authMW)
	g.POST("/", h.Chat)
}

// Chat godoc
// @Summary Отправить сообщение в LLM
// @Tags Chat
// @Description Отправляет сообщение пользователя в языковую модель и возвращает ответ
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param input body ChatRequest true "Входное сообщение"
// @Security BearerAuth
// @Success 200 {object} ChatResponse "Ответ от LLM"
// @Failure 400 {object} gin.H "Неверный запрос"
// @Failure 401 {object} gin.H "Неавторизованный доступ"
// @Failure 500 {object} gin.H "Внутренняя ошибка сервера"
// @Router /chat [post]
func (h *Handler) Chat(c *gin.Context) {
	userID, ok := c.MustGet("userID").(int64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный идентификатор пользователя"})
		return
	}

	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	messages, err := h.service.Chat(c.Request.Context(), userID, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, ChatResponse{Messages: messages})
}
