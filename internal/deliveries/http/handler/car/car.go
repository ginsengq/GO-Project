package carhandler

import (
	"errors"
	"net/http"
	"strconv"

	"myproject/internal/entities"
	"myproject/internal/usecases/car"
	"myproject/pkg/logger"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc     car.CarUseCase
	logger logger.Interface
}

func NewHandler(uc car.CarUseCase, logger logger.Interface) *Handler {
	return &Handler{uc: uc, logger: logger}
}

func (h *Handler) CreateCar(c *gin.Context) {
	var input entities.Car
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	createdCar, err := h.uc.CreateCar(c.Request.Context(), &input)
	if err != nil {
		h.logger.Error("create car failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, createdCar)
}

func (h *Handler) GetCar(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	car, err := h.uc.GetCar(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "car not found"})
			return
		}
		h.logger.Error("get car failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, car)
}

func (h *Handler) UpdateCar(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var input entities.CarUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	updatedCar, err := h.uc.UpdateCar(c.Request.Context(), id, input)
	if err != nil {
		h.logger.Error("update car failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, updatedCar)
}

func (h *Handler) DeleteCar(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.uc.DeleteCar(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("delete car failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) ListCars(c *gin.Context) {
	var filter entities.CarFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}

	cars, total, err := h.uc.ListCars(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error("list cars failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": cars, "total": total})
}

func (h *Handler) ChangeCarStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var input struct {
		Status entities.CarStatus `json:"status"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	car, err := h.uc.ChangeCarStatus(c.Request.Context(), id, input.Status)
	if err != nil {
		h.logger.Error("change car status failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, car)
}
