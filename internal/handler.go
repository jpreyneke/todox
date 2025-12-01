package internal

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
		v1.POST("/todos", h.CreateTodos)
		v1.PATCH("/todos", h.UpdateTodos)
		v1.GET("/todos", h.ListTodos)
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) CreateTodos(c *gin.Context) {
	var body struct {
		Todos []CreateTodoInput `json:"todos" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	todos, err := h.service.BulkCreate(c.Request.Context(), body.Todos)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": todos})
}

func (h *Handler) UpdateTodos(c *gin.Context) {
	var body struct {
		Todos []UpdateTodoInput `json:"todos" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	todos, err := h.service.BulkUpdate(c.Request.Context(), body.Todos)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": todos})
}

func (h *Handler) ListTodos(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'page' parameter"})
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'limit' parameter"})
		return
	}

	todos, total, err := h.service.List(c.Request.Context(), page, limit)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": todos,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		c.JSON(http.StatusNotFound, ErrorResponse{Error: ErrNotFound.Error()})
	case errors.Is(err, ErrDuplicateTitle):
		c.JSON(http.StatusConflict, ErrorResponse{Error: ErrDuplicateTitle.Error()})
	case errors.Is(err, ErrTitleRequired):
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: ErrTitleRequired.Error()})
	case errors.Is(err, ErrTitleEmpty):
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: ErrTitleEmpty.Error()})
	case errors.Is(err, ErrTitleMaxLength):
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: ErrTitleMaxLength.Error()})
	case errors.Is(err, ErrInvalidID):
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: ErrInvalidID.Error()})
	case errors.Is(err, ErrEmptyList):
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: ErrEmptyList.Error()})
	case errors.Is(err, ErrDuplicateInRequest):
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: ErrDuplicateInRequest.Error()})
	case errors.Is(err, ErrLimitExceeded):
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: ErrLimitExceeded.Error()})
	default:
		// Log unexpected errors
		slog.Error("Unexpected error",
			"error", err,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP(),
		)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
	}
}
