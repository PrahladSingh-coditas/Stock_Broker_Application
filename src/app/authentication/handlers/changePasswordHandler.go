package handlers

import (
	"authentication/business"
	"authentication/commons/constants"
	"authentication/models"

	"github.com/gin-gonic/gin"
)

type ChangePasswordHandler struct {
	service *business.ChangePasswordService
}

func NewChangePasswordHandler(service *business.ChangePasswordService) *ChangePasswordHandler {
	return &ChangePasswordHandler{
		service: service,
	}
}

// HandleChangePassword godoc
// @Summary Reset user password
// @Description Reset password
// @Tags password-reset-api
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Security BearerAuth
// @Param request body models.BFFResetPasswordRequest true "Reset Password Request"
// @Success 200 {object} models.ErrorAPIResponse "password reset successful"
// @Failure 400 {object} models.ErrorAPIResponse "bad request"
// @Failure 401 {object} models.ErrorAPIResponse "unauthorized"
// @Failure 401 {object} models.ErrorAPIResponse "passwords miss match"
// @Failure 500 {object} models.ErrorAPIResponse "internal server error"
// @Router /api/auth/reset-password [post]
func (controller *ChangePasswordHandler) HandleChangePassword(c *gin.Context) {
	var req models.BFFResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(400, gin.H{"error": "Invalid request payload"})
		return
	}
	if req.Password != req.ConfirmPassword {
		c.IndentedJSON(400, gin.H{"error": "Passwords miss match"})
		return
	}

	//get username from middleware
	username, exists := c.Get("username")
	usernamestr := username.(string)
	if !exists {
		c.IndentedJSON(400, gin.H{"error": "Username not found"})
		return
	}

	err := controller.service.UpdatePassword(usernamestr, req)
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(200, gin.H{
		"message": constants.PasswordChangeSuccessMsg,
	})

}
