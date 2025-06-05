package user

import (
	"errors"
	"net/http"

	userSvc "github.com/board-box/backend/internal/service/user"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *userSvc.Service
	authMW  func(c *gin.Context)
}

func New(service *userSvc.Service, authMW func(c *gin.Context)) *Handler {
	return &Handler{service: service, authMW: authMW}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/user")
	g.POST("/register", h.Register)
	g.POST("/login", h.Login)

	g.Use(h.authMW)
	g.GET("/info", h.Info)
}

// Register godoc
// @Summary Регистрация пользователя
// @Tags Users
// @Description Регистрирует нового пользователя с email и паролем
// @Accept json
// @Produce json
// @Param input body RegisterRequest true "Данные для регистрации"
// @Success 201
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /user/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	err := h.service.Register(c.Request.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save user"})
		return
	}

	c.Status(http.StatusCreated)
}

// Login godoc
// @Summary Авторизация пользователя
// @Tags Users
// @Description Авторизует пользователя и возвращает JWT токен
// @Accept json
// @Produce json
// @Param input body LoginRequest true "Данные для входа"
// @Success 200 {object} map[string]string "token"
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /user/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	token, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, userSvc.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Info godoc
// @Summary Получить информацию о пользователе
// @Tags Users
// @Description Возвращает информацию о текущем авторизованном пользователе
// @Security BearerAuth
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Success 200 {object} InfoResponse "Информация о пользователе"
// @Failure 401 {object} gin.H "Неавторизованный доступ"
// @Failure 500 {object} gin.H "Внутренняя ошибка сервера"
// @Router /user/info [get]
func (h *Handler) Info(c *gin.Context) {
	userID, ok := c.MustGet("userID").(int64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный идентификатор пользователя"})
		return
	}

	info, err := h.service.Info(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, userSvc.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении информации"})
		return
	}

	c.JSON(http.StatusOK, InfoResponse{
		Username: info.Username,
		Email:    info.Email,
	})
}
