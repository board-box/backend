package collection

import (
	"errors"
	"net/http"
	"strconv"

	collectionSvc "github.com/board-box/backend/internal/service/collection"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *collectionSvc.Service
	authMW  func(c *gin.Context)
}

func New(service *collectionSvc.Service, authMW func(c *gin.Context)) *Handler {
	return &Handler{service: service, authMW: authMW}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/collections")

	g.Use(h.authMW)

	g.GET("/", h.ListCollections)
	g.GET("/:id", h.GetCollection)
	g.POST("/", h.CreateCollection)
	g.PUT("/:id", h.UpdateCollection)
	g.DELETE("/:id", h.DeleteCollection)
	g.POST("/:id/games/:game_id", h.AddGameToCollection)
	g.DELETE("/:id/games/:game_id", h.RemoveGameFromCollection)
}

// ListCollections godoc
// @Summary Список коллекций пользователя
// @Tags Collections
// @Description Получить список всех коллекций текущего пользователя
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Security BearerAuth
// @Success 200 {array} collectionSvc.Collection
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /collections [get]
func (h *Handler) ListCollections(c *gin.Context) {
	userID, ok := c.MustGet("userID").(int64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	collections, err := h.service.ListCollections(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить коллекции"})
		return
	}
	c.JSON(http.StatusOK, collections)
}

// GetCollection godoc
// @Summary Получить коллекцию по ID
// @Tags Collections
// @Description Получить коллекцию по её идентификатору
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param id path int true "ID коллекции"
// @Security BearerAuth
// @Success 200 {object} collectionSvc.Collection
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /collections/{id} [get]
func (h *Handler) GetCollection(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID"})
		return
	}

	userID, ok := c.MustGet("userID").(int64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	collection, err := h.service.GetCollection(c.Request.Context(), id, userID)
	if err != nil {
		if errors.Is(err, collectionSvc.ErrCollectionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Коллекция не найдена"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить коллекцию"})
		return
	}

	c.JSON(http.StatusOK, collection)
}

// CreateCollection godoc
// @Summary Создать новую коллекцию
// @Tags Collections
// @Description Создать новую коллекцию для текущего пользователя
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param input body CreateCollectionRequest true "Данные коллекции"
// @Security BearerAuth
// @Success 201 {object} collectionSvc.Collection
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /collections [post]
func (h *Handler) CreateCollection(c *gin.Context) {
	userID, ok := c.MustGet("userID").(int64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	var req CreateCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	collection, err := h.service.CreateCollection(c.Request.Context(), userID, convertCreateReqToDTO(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать коллекцию"})
		return
	}

	c.JSON(http.StatusCreated, collection)
}

// UpdateCollection godoc
// @Summary Обновить коллекцию
// @Tags Collections
// @Description Обновить данные коллекции
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param id path int true "ID коллекции"
// @Param input body UpdateCollectionRequest true "Новые данные коллекции"
// @Security BearerAuth
// @Success 200 {object} collectionSvc.Collection
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /collections/{id} [put]
func (h *Handler) UpdateCollection(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID"})
		return
	}

	userID, ok := c.MustGet("userID").(int64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	var req UpdateCollectionRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	collection, err := h.service.UpdateCollection(c.Request.Context(), id, userID, convertUpdateReqToDTO(req))
	if err != nil {
		if errors.Is(err, collectionSvc.ErrCollectionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Коллекция не найдена"})
			return
		}
		if errors.Is(err, collectionSvc.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Нет прав для изменения коллекции"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить коллекцию"})
		return
	}

	c.JSON(http.StatusOK, collection)
}

// DeleteCollection godoc
// @Summary Удалить коллекцию
// @Tags Collections
// @Description Удалить коллекцию по её ID
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param id path int true "ID коллекции"
// @Security BearerAuth
// @Success 204
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /collections/{id} [delete]
func (h *Handler) DeleteCollection(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID"})
		return
	}

	userID, ok := c.MustGet("userID").(int64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.service.DeleteCollection(c.Request.Context(), id, userID); err != nil {
		if errors.Is(err, collectionSvc.ErrCollectionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Коллекция не найдена"})
			return
		}
		if errors.Is(err, collectionSvc.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Нет прав для удаления коллекции"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить коллекцию"})
		return
	}

	c.Status(http.StatusNoContent)
}

// AddGameToCollection godoc
// @Summary Добавить игру в коллекцию
// @Tags Collections
// @Description Добавляет игру в коллекцию пользователя
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param id path int true "ID коллекции"
// @Param game_id path int true "ID игры"
// @Security BearerAuth
// @Success 204
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /collections/{id}/games/{game_id} [post]
func (h *Handler) AddGameToCollection(c *gin.Context) {
	collectionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат id"})
		return
	}

	gameID, err := strconv.ParseInt(c.Param("game_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат game_id"})
		return
	}

	userID, ok := c.MustGet("userID").(int64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.service.AddGameToCollection(c.Request.Context(), collectionID, gameID, userID)
	if err != nil {
		if errors.Is(err, collectionSvc.ErrCollectionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Коллекция не найдена"})
			return
		}
		if errors.Is(err, collectionSvc.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Нет прав на изменение коллекции"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось добавить игру в коллекцию"})
		return
	}

	c.Status(http.StatusNoContent)
}

// RemoveGameFromCollection godoc
// @Summary Удалить игру из коллекции
// @Tags Collections
// @Description Удаляет игру из коллекции пользователя
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param id path int true "ID коллекции"
// @Param game_id path int true "ID игры"
// @Security BearerAuth
// @Success 204
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /collections/{id}/games/{game_id} [delete]
func (h *Handler) RemoveGameFromCollection(c *gin.Context) {
	collectionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат id"})
		return
	}

	gameID, err := strconv.ParseInt(c.Param("game_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат game_id"})
		return
	}

	userID, ok := c.MustGet("userID").(int64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.service.RemoveGameFromCollection(c.Request.Context(), collectionID, gameID, userID)
	if err != nil {
		if errors.Is(err, collectionSvc.ErrCollectionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Коллекция не найдена"})
			return
		}
		if errors.Is(err, collectionSvc.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Нет прав на изменение коллекции"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить игру из коллекции"})
		return
	}

	c.Status(http.StatusNoContent)
}
