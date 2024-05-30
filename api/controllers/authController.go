package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/berkaycubuk/billiard_software_api/database"
	"github.com/berkaycubuk/billiard_software_api/models"
	"github.com/berkaycubuk/billiard_software_api/pkg/brevo"
	"github.com/berkaycubuk/billiard_software_api/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required" validate:"required"`
	NewPassword string `json:"new_password" binding:"required" validate:"required"`
}

func UpdatePassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request UpdatePasswordRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		// get user
		user, err := utils.CurrentUser(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.user_not_found",
				"success": false,
			})
			return
		}

		// check current password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.CurrentPassword))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.current_password_wrong",
				"success": false,
			})
			return
		}

		// update password

		// hash password
		hashed, _ := bcrypt.GenerateFromPassword([]byte(request.NewPassword), 8)

		err = database.DB.Model(&user).Update("password", string(hashed)).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.unable_to_update_password",
				"success": false,
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		})
	}
}

type PasswordResetCompleteRequest struct {
	Token		string	`json:"token" binding:"required" validate:"required"`
	Password		string	`json:"password" binding:"required" validate:"required"`
}

func PasswordResetComplete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request PasswordResetCompleteRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		var passwordReset models.PasswordReset
		err := database.DB.Where("token = ? AND valid = ? AND expires_at > ?", request.Token, true, time.Now()).First(&passwordReset).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"success": false,
			})
			return
		}

		// get user
		var user models.User
		err = database.DB.Where("id = ?", passwordReset.UserID).First(&user).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"success": false,
			})
			return
		}

		// hash password
		hashed, _ := bcrypt.GenerateFromPassword([]byte(request.Password), 8)

		err = database.DB.Model(&user).Update("password", string(hashed)).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.unable_to_update_password",
				"success": false,
			})
			return
		}

		// expire the password reset token
		database.DB.Model(&passwordReset).Update("valid", false)

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

type VerifyRequest struct {
	Token		string	`json:"token" binding:"required" validate:"required"`
}

func Verify() gin.HandlerFunc {
	return func (ctx *gin.Context)  {
		var request VerifyRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		// find user
		var user models.User
		err := database.DB.Where("verify_token = ? AND email_verified_at IS NULL", request.Token).First(&user).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.user_not_found",
				"success": false,
			})
			return
		}

		database.DB.Model(&user).Update("email_verified_at", time.Now())
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

type ForgotPasswordRequest struct {
	Email		string	`json:"email" binding:"required" validate:"required"`
}

func ForgotPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request ForgotPasswordRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		// find user
		var user models.User
		err := database.DB.Where("email = ?", request.Email).First(&user).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.user_not_found",
				"success": false,
			})
			return
		}

		// create password reset row and token
		token := utils.RandomString(10)

		database.DB.Create(&models.PasswordReset{
			UserID: user.ID,
			Token: token,
			Valid: true,
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(time.Hour),
		})
		passwordResetLink := os.Getenv("FRONTEND_URL") + "/reset-password-continue/" + token

		// send mail
		err = brevo.SendEmail(user.Email, user.Name, "Password Reset", "You requested a password reset.<br/><br/>Here is your link to reset your password: <a href='" + passwordResetLink + "'>Reset password</a><br/><br/>If you don't requested it, there is no need to click the link and you can forget this mail.")

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"success": false,
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "success.mail_sent",
			"success": true,
		})
	}
}

type LoginRequest struct {
	Email		string	`json:"email" binding:"required" validate:"required"`
	Password	string	`json:"password" binding:"required" validate:"required"`
}

func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var loginRequest LoginRequest

		if !utils.ValidateRequest(ctx, &loginRequest) {
			return
		}

		fmt.Println(loginRequest.Email)

		var user models.User
		dbs := database.DB.Where("email = ?", loginRequest.Email).First(&user)
		if dbs.Error != nil {
			fmt.Println("--- HIT ---")
			fmt.Println(loginRequest.Email)
			fmt.Println("--- HIT ---")
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.user_not_found",
				"success": false,
			})
			return
		}

		// check passwords matching
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.user_not_found",
				"success": false,
			})
			return
		}

		token, err := utils.GenerateToken(&user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.unable_to_generate_token",
				"success": false,
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"token": token,
			"success": true,
		})
	}
}

type RegisterRequest struct {
	Name		string	`json:"name" binding:"required" validate:"required"`
	Surname		string	`json:"surname" binding:"required" validate:"required"`
	Phone		string	`json:"phone" binding:"required" validate:"required"`
	Email		string	`json:"email" binding:"required,email" validate:"required"`
	Password	string	`json:"password" binding:"required" validate:"required"`
	ConfirmPassword	string	`json:"confirm_password" binding:"required" validate:"required"`
}

func Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var registerRequest RegisterRequest

		if !utils.ValidateRequest(ctx, &registerRequest) {
			return
		}

		// check passwords matching
		if registerRequest.Password != registerRequest.ConfirmPassword {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.passwords_do_not_match",
				"success": false,
			})
			return
		}

		// hash password
		hashed, _ := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), 8)

		// create user model
		token := utils.RandomString(10)
		user := models.User{
			Name: registerRequest.Name,
			Surname: registerRequest.Surname,
			Phone: registerRequest.Phone,
			Email: registerRequest.Email,
			VerifyToken: token,
			Password: string(hashed),
		}	

		// create user row
		dbc := database.DB.Create(&user)
		if dbc.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.email_or_phone_already_in_use",
				"success": false,
			})
			return
		}

		// create user role
		database.DB.Create(&models.UserRole{
			UserID: user.ID,
			RoleID: 2, // user
		})

		// confirmation link
		confirmationLink := os.Getenv("FRONTEND_URL") + "/verify/" + token

		// send mail
		brevo.SendEmail(user.Email, user.Name, "Confirm your account", "An account created with your e-mail.<br/><br/>If this happened under your permission please confirm your account: <a href='" + confirmationLink + "'>Confirm account</a><br/><br/>If you don't requested it, there is no need to click the link and you can forget this mail.")

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}
