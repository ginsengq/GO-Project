package handler

import (
	"myproject/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommonHandler struct {
	logger logger.Interface
}

func NewCommonHandler(logger logger.Interface) *CommonHandler {
	return &CommonHandler{logger: logger}
}

func (h *CommonHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "available",
		"version": "1.0.0",
	})
}

func (h *CommonHandler) NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"error":   "not found",
		"message": "requested resource was not found",
	})
}

func (h *CommonHandler) MethodNotAllowed(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{
		"error":   "method not allowed",
		"message": "this http method is not supported for the requested resource",
	})
}
