package orderhandler

import (
	"net/http"
	"strconv"

	"myproject/internal/entities"
	ordercase "myproject/internal/usecases/order"
	"myproject/pkg/logger"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	orderUC ordercase.UseCase
	logger  logger.Interface
}

func NewHandler(orderUC ordercase.UseCase, logger logger.Interface) *Handler {
	return &Handler{orderUC: orderUC, logger: logger}
}

type CreateOrderRequest struct {
	UserID     int     `json:"user_id" binding:"required,gt=0"`
	CarID      int     `json:"car_id" binding:"required,gt=0"`
	Deposit    float64 `json:"deposit" binding:"gte=0"`
	TotalPrice float64 `json:"total_price" binding:"required,gt=0"`
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("CreateOrder: invalid input", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	order := &entities.Order{
		UserID:     req.UserID,
		CarID:      req.CarID,
		Deposit:    req.Deposit,
		TotalPrice: req.TotalPrice,
	}

	orderID, err := h.orderUC.CreateOrder(c.Request.Context(), order)
	if err != nil {
		h.logger.Error("CreateOrder: failed to create order", "order", order, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": orderID, "message": "order created successfully"})
}

func (h *Handler) GetOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		h.logger.Error("GetOrder: invalid id", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	order, err := h.orderUC.GetOrder(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("GetOrder: failed to get order", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *Handler) GetOrdersByUserID(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		h.logger.Error("GetOrdersByUserID: invalid user_id", "user_id", userIDStr, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	orders, err := h.orderUC.GetOrdersByUserID(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("GetOrdersByUserID: failed to get orders", "user_id", userID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		h.logger.Error("UpdateOrderStatus: invalid id", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("UpdateOrderStatus: invalid input", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	err = h.orderUC.UpdateOrderStatus(c.Request.Context(), id, req.Status)
	if err != nil {
		h.logger.Error("UpdateOrderStatus: failed to update status", "id", id, "status", req.Status, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order status updated successfully"})
}

func (h *Handler) CancelOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		h.logger.Error("CancelOrder: invalid id", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.orderUC.CancelOrder(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("CancelOrder: failed to cancel order", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order cancelled successfully"})
}

func (h *Handler) ListAllOrders(c *gin.Context) {
	orders, err := h.orderUC.ListAllOrders(c.Request.Context())
	if err != nil {
		h.logger.Error("ListAllOrders: failed to list all orders", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}
