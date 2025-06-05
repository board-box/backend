package handler

import (
	"errors"
	"net/http"
	"strconv"

	gameSvc "github.com/board-box/backend/internal/service/game"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *gameSvc.Service
}

func New(service *gameSvc.Service) *Handler {
	return &Handler{service}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/games")
	g.GET("/", h.ListGames)
	g.GET("/:id", h.GetGame)
	g.POST("/by-ids", h.GetGamesByIDs)
}

// ListGames godoc
// @Summary Список игр
// @Tags Games
// @Description Получить список всех настольных игр
// @Produce json
// @Success 200 {array} gameSvc.Game
// @Failure 500 {object} gin.H
// @Router /games/ [get]
func (h *Handler) ListGames(c *gin.Context) {
	games, err := h.service.ListGames(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить игры"})
		return
	}
	c.JSON(http.StatusOK, games)
}

// GetGame godoc
// @Summary Получить игру по ID
// @Tags Games
// @Description Получить настольную игру по её идентификатору
// @Produce json
// @Param id path string true "ID игры"
// @Success 200 {object} gameSvc.Game
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /games/{id} [get]
func (h *Handler) GetGame(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID игры обязателен"})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID должен быть int"})
		return
	}

	game, err := h.service.GetGame(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gameSvc.ErrGameNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Игра не найдена"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить игру"})
		return
	}

	c.JSON(http.StatusOK, game)
}

// GetGamesByIDs godoc
// @Summary Получить список игр по ID
// @Tags Games
// @Description Получить несколько игр по их идентификаторам
// @Accept json
// @Produce json
// @Param input body GetGamesByIDsRequest true "Список ID игр"
// @Success 200 {array} gameSvc.Game
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /games/by-ids [post]
func (h *Handler) GetGamesByIDs(c *gin.Context) {
	var req GetGamesByIDsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса или пустой список ID"})
		return
	}

	games, err := h.service.GetGames(c.Request.Context(), req.IDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить игры"})
		return
	}

	c.JSON(http.StatusOK, games)
}
