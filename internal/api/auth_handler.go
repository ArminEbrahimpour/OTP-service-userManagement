package api

import (
	"net/http"
	"otp-service/internal/models"
	"otp-service/internal/service"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// sendOTP godoc
// @Summary Send otp to phone number
// @Description Send a 6 digit OTP code to the specific phone number
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.OTPRequest true "Phone number"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 429 {object} map[string]string
// @Router /api/v1/auth/send-otp [post]
func (h *AuthHandler) SendOTP(c *gin.Context) {
	var req models.OTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.SendOTP(req.PhoneNumber); err != nil {
		if err.Error() == "rate limit exceeded" || strings.Contains(err.Error(), "rate limit exceeded") {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

// verifyOTP godoc
// @Summary Verify OTP and authenticate user
// @Description Verify the OTP code and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.OTPVerifyReq true "Phone number and OTP code"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 429 {object} map[string]string
// @Router /auth/verify-otp [post]
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req models.OTPVerifyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := h.authService.VerifyOTP(req.PhoneNumber, req.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}
