package userhandler

import (
	"errors"
	"net/http"
	"strconv"

	"myproject/internal/entities"
	usercase "myproject/internal/usecases/user"

	"myproject/pkg/logger"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	userUC usercase.UseCase
	logger logger.Interface
}

func NewHandler(userUC usercase.UseCase, logger logger.Interface) *Handler {
	return &Handler{userUC: userUC, logger: logger}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var input entities.User
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Error("CreateUser: invalid input", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := h.userUC.Create(c.Request.Context(), &input); err != nil {
		h.logger.Error("CreateUser: user creation failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})
}

func (h *Handler) GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		h.logger.Error("GetUserByID: invalid id", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	user, err := h.userUC.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, entities.ErrInvalidID) {
			h.logger.Warn("GetUserByID: invalid user ID", "id", id)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		if errors.Is(err, entities.ErrNotFound) {
			h.logger.Warn("GetUserByID: user not found", "id", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		h.logger.Error("GetUserByID: failed to get user", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		h.logger.Error("UpdateUser: invalid id", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var input entities.User
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Error("UpdateUser: invalid input", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	input.ID = id

	if err := h.userUC.Update(c.Request.Context(), &input); err != nil {
		if errors.Is(err, entities.ErrInvalidID) {
			h.logger.Warn("UpdateUser: invalid user ID", "id", id)
			c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required in the request body"})
			return
		}
		if errors.Is(err, entities.ErrNotFound) {
			h.logger.Warn("UpdateUser: user not found", "id", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		h.logger.Error("UpdateUser: failed to update user", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		h.logger.Error("DeleteUser: invalid id", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.userUC.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, entities.ErrInvalidID) {
			h.logger.Warn("DeleteUser: invalid user ID", "id", id)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		h.logger.Error("DeleteUser: failed to delete user", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) ListUsers(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		h.logger.Error("ListUsers: invalid limit", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		h.logger.Error("ListUsers: invalid offset", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
		return
	}

	users, err := h.userUC.List(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("ListUsers: failed to list users", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	count, err := h.userUC.Count(c.Request.Context())
	if err != nil {
		h.logger.Error("ListUsers: failed to count users", "error", err)
	}

	c.JSON(http.StatusOK, gin.H{"items": users, "total": count})
}

func (h *Handler) ChangePassword(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		h.logger.Error("ChangePassword: invalid id", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var input struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Error("ChangePassword: invalid input", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	err = h.userUC.ChangePassword(c.Request.Context(), id, input.OldPassword, input.NewPassword)
	if err != nil {
		if errors.Is(err, entities.ErrInvalidID) { // Используем ошибки из entity
			h.logger.Warn("ChangePassword: invalid user ID", "id", id)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			return
		}
		if errors.Is(err, errors.New("new password is required")) {
			h.logger.Warn("ChangePassword: new password required", "id", id)
			c.JSON(http.StatusBadRequest, gin.H{"error": "new password is required"})
			return
		}
		if errors.Is(err, errors.New("invalid old password")) {
			h.logger.Warn("ChangePassword: invalid old password", "id", id)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid old password"})
			return
		}
		h.logger.Error("ChangePassword: failed to change password", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}

func (h *Handler) AuthenticateUser(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Error("AuthenticateUser: invalid input", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	token, user, err := h.userUC.Authenticate(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		if errors.Is(err, errors.New("email and password are required")) {
			h.logger.Warn("AuthenticateUser: email and password required")
			c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
			return
		}
		if errors.Is(err, entities.ErrNotFound) {
			h.logger.Warn("AuthenticateUser: invalid credentials (user not found)")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		if errors.Is(err, errors.New("invalid credentials")) {
			h.logger.Warn("AuthenticateUser: invalid credentials")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		h.logger.Error("AuthenticateUser: authentication failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}
