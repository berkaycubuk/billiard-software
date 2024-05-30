package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/berkaycubuk/billiard_software_api/database"
	"github.com/berkaycubuk/billiard_software_api/models"
	"github.com/berkaycubuk/billiard_software_api/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetGuests() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var gameUsers []models.GameUser
		err := database.DB.Where("user_id IS NULL AND ended_at IS NULL").Find(&gameUsers).Error
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"guests":  nil,
				"success": true,
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"guests":  &gameUsers,
			"success": true,
		})
	}
}

func GetUsersWithoutPagination() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var users []models.User
		err := database.DB.Where("deleted_at IS NULL").Order("id desc").Find(&users).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error.users_not_found",
				"success": false,
			})
			return
		}

		ctx.Set("response", gin.H{
			"users": users,
		})
	}
}

func GetUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userParam := ctx.Query("user")

		var usersTemp []models.User
		usersQuery := database.DB

		if userParam != "" {
			//usersQuery = usersQuery.Where("id = ?", userParam)
			usersQuery = usersQuery.Where("name LIKE ?", "%" + userParam + "%")
		}

		usersQuery.Where("deleted_at IS NULL").Find(&usersTemp)
		count := len(usersTemp)

		pageSize := count / 10
		if pageSize*10 < count {
			pageSize = pageSize + 1
		}

		var users []models.User
		realUsersQuery := database.DB

		if userParam != "" {
			//realUsersQuery = realUsersQuery.Where("id = ?", userParam)
			realUsersQuery = realUsersQuery.Where("name LIKE ?", "%" + userParam + "%")
		}

		err := realUsersQuery.Where("deleted_at IS NULL").Scopes(utils.Paginate(ctx)).Order("id desc").Find(&users).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error.users_not_found",
				"success": false,
			})
			return
		}

		ctx.Set("response", gin.H{
			"users": users,
			"pagination": gin.H{
				"size": pageSize,
			},
		})
	}
}

type GetTableRequest struct {
	ID string `uri:"id" binding:"required"`
}

func GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request GetTableRequest
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "ID required.",
				"success": false,
			})
			return
		}

		var user models.User
		err := database.DB.Where("id = ?", request.ID).Preload("Subscription").Preload("Role.Role").First(&user).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.user_not_found",
				"success": false,
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"user":    &user,
			"success": true,
		})
	}
}

type GetMyActiveSubscriptionResponse struct {
	Subscription *models.UserSubscription `json:"subscription"`
	Success      bool                     `json:"success"`
}

func GetMyActiveSubscription() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := utils.CurrentUser(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error.user_not_found",
				"success": false,
			})
			return
		}

		var subscription models.UserSubscription
		err = database.DB.Where("user_id = ? AND status IN (?, ?)", user.ID, models.USER_SUBSCRIPTION_STATUS_ACTIVE, models.USER_SUBSCRIPTION_STATUS_PAUSED).First(&subscription).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.JSON(http.StatusOK, GetMyActiveSubscriptionResponse{
					Subscription: nil,
					Success:      true,
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"success": false,
			})
			return
		}

		ctx.JSON(http.StatusOK, GetMyActiveSubscriptionResponse{
			Subscription: &subscription,
			Success:      true,
		})
	}
}

func GetMySubscriptions() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := utils.CurrentUser(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error.user_not_found",
				"success": false,
			})
			return
		}

		var subscriptions []models.UserSubscription
		err = database.DB.Where("user_id = ?", user.ID).Find(&subscriptions).Error
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"success": false,
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"subscriptions": subscriptions,
			"success":       true,
		})
	}
}

type ProfileResponse struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	Surname    string    `json:"surname"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	RoleID     uint      `json:"role_id"`
	IsVerified bool      `json:"is_verified"`
	IsAdmin    bool      `json:"is_admin"`
	CreatedAt  time.Time `json:"created_at"`
}

func GetProfile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := utils.CurrentUser(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error.user_not_found",
				"success": false,
			})
			return
		}

		/*
		var userRole models.UserRole
		err = database.DB.Where("user_id = ?", user.ID).First(&userRole).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error.user_role_not_found",
				"success": false,
			})
			return
		}
		*/

		ctx.JSON(http.StatusOK, gin.H{
			"user": ProfileResponse{
				ID:         user.ID,
				Name:       user.Name,
				Surname:    user.Surname,
				Email:      user.Email,
				Phone:      user.Phone,
				RoleID:     user.Role.RoleID,
				IsVerified: user.EmailVerifiedAt != nil,
				IsAdmin:    user.Role.RoleID == 1,
				CreatedAt:  user.CreatedAt,
			},
			"success": true,
		})
	}
}

type GetPermissionsPermission struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}

type GetPermissionsResponse struct {
	Permissions *[]GetPermissionsPermission `json:"permissions"`
	Success     bool                        `json:"success"`
}

func permissionValueFormatter(valueType string, value string) any {
	if valueType == "bool" {
		return value == "true"
	}

	return value
}

func GetPermissions() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := utils.CurrentUser(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error.user_not_found",
				"success": false,
			})
			return
		}

		var userRole models.UserRole
		err = database.DB.Where("user_id = ?", user.ID).First(&userRole).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error.user_role_not_found",
				"success": false,
			})
			return
		}

		var permissions []GetPermissionsPermission
		var roleConfigs []models.RoleConfig
		err = database.DB.Where("role_id = ?", userRole.RoleID).Find(&roleConfigs).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error.role_configs_not_found",
				"success": false,
			})
			return
		}

		for _, v := range roleConfigs {
			permissions = append(permissions, GetPermissionsPermission{
				Name:  v.Name,
				Value: permissionValueFormatter(v.ValueType, v.Value),
			})
		}

		ctx.JSON(http.StatusOK, GetPermissionsResponse{
			Permissions: &permissions,
			Success:     true,
		})
	}
}

type CreateUserRequest struct {
	Name    string `json:"name" binding:"required" validate:"required"`
	Surname string `json:"surname" binding:"required" validate:"required"`
	Email   string `json:"email" binding:"required" validate:"required"`
	Phone   string `json:"phone" binding:"required" validate:"required"`
	//EmailVerified bool   `json:"email_verified" binding:"required" validate:"required"`
	Role   uint `json:"role"`
}

func CreateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request CreateUserRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		user := models.User{
			Name:     request.Name,
			Surname:  request.Surname,
			Email:    request.Email,
			Phone:    request.Phone,
			Password: "$2a$08$pUpPKofXRhkCp2bslK91veosI3F2v8tMzOMY2gFnYKb19MHUdxILS", // password
		}

		err := database.DB.Create(&user).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"success": false,
			})
		}

		database.DB.Create(&models.UserRole{
				UserID: user.ID,
				RoleID: request.Role,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
		})

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

type UpdateUserRequest struct {
	ID      uint   `json:"id" binding:"required" validate:"required"`
	Name    string `json:"name" binding:"required" validate:"required"`
	Surname string `json:"surname" binding:"required" validate:"required"`
	Email   string `json:"email" binding:"required" validate:"required"`
	Phone   string `json:"phone" binding:"required" validate:"required"`
	Role   uint `json:"role"`
}

func UpdateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request UpdateUserRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		// check user exists
		var user models.User
		err := database.DB.Where("id = ?", request.ID).First(&user).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.user_not_found",
				"success": false,
			})
		}

		err = database.DB.Model(&user).Updates(models.User{
			Name:    request.Name,
			Surname: request.Surname,
			Email:   request.Email,
			Phone:   request.Phone,
		}).Error
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"success": false,
			})
		}

		var userRole models.UserRole
		err = database.DB.Where("user_id = ?", request.ID).First(&userRole).Error
		if err == nil {
				database.DB.Model(&userRole).Update("role_id", request.Role)
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

type DeleteUserRequest struct {
	ID uint `uri:"id" binding:"required" validate:"required"`
}

func DeleteUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request DeleteUserRequest
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.id_required",
				"success": false,
			})
			return
		}

		// check user exists
		var user models.User
		err := database.DB.Where("id = ?", request.ID).First(&user).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.user_not_found",
				"success": false,
			})
		}

		err = database.DB.Model(&user).Update("deleted_at", time.Now()).Error
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"success": false,
			})
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}
