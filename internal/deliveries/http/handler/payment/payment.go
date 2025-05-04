package paymenthandler

import (
	"net/http"
	"strconv"

	"myproject/internal/entities"
	paymentcase "myproject/internal/usecases/payment"
	"myproject/pkg/logger"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	paymentUC paymentcase.PaymentUseCase
	logger    logger.Interface
}

func NewHandler(paymentUC paymentcase.PaymentUseCase, logger logger.Interface) *Handler {
	return &Handler{paymentUC: paymentUC, logger: logger}
}

type DepositRequest struct {
	UserID int     `json:"user_id" binding:"required,gt=0"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

func (h *Handler) Deposit(c *gin.Context) {
	var req DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Deposit: invalid input", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := h.paymentUC.Deposit(c.Request.Context(), req.UserID, req.Amount); err != nil {
		h.logger.Error("Deposit: failed to deposit", "user_id", req.UserID, "amount", req.Amount, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deposit successful"})
}

type CreateTransactionRequest struct {
	UserID int     `json:"user_id" binding:"required,gt=0"`
	Amount float64 `json:"amount" binding:"required"`
	Type   string  `json:"type" binding:"required"`
}

func (h *Handler) CreateTransaction(c *gin.Context) {
	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("CreateTransaction: invalid input", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	tx := &entities.Transaction{
		UserID: req.UserID,
		Amount: req.Amount,
		Type:   req.Type,
	}

	if err := h.paymentUC.CreateTransaction(c.Request.Context(), tx); err != nil {
		h.logger.Error("CreateTransaction: failed to create transaction", "transaction", tx, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "transaction created"})
}

func (h *Handler) GetTransactionsByUser(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		h.logger.Error("GetTransactionsByUser: invalid user_id", "user_id", userIDStr, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	transactions, err := h.paymentUC.GetTransactionsByUser(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("GetTransactionsByUser: failed to get transactions", "user_id", userID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}
