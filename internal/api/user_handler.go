package api

import (
	"net/http"
	"otp-service/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetUser godoc
// @Summary Get user by ID
// @Description get a single user by speicified ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 429 {object} map[string]string
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user)

}

// GetUsers godoc
// @Summary Get list of users
// @Description Get a paginated users list with optional search
// @Tags users
// @Accept json
// @Produce json
// @Param page_size query int false "Page size" default(10)
// @Param page query int false "Page number" default(1)
// @Param search query string false "Search by phone number"
// @Success 200 {object} models.UserListResponse
// @Failure 400 {object} map[string]string
// @Failure 429 {object} map[string]string
// @Router /api/v1/users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	page := 1
	pageSize := 10
	search := c.Query("search")

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if sizeStr := c.Query("page_size"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 && s <= 100 {
			pageSize = s
		}
	}

	response, err := h.userService.GetUsers(page, pageSize, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response)
}
