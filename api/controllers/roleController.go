package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/berkaycubuk/billiard_software_api/database"
	"github.com/berkaycubuk/billiard_software_api/models"
	"github.com/berkaycubuk/billiard_software_api/utils"
	"github.com/gin-gonic/gin"
)

type GetRoleRequest struct {
	ID uint `uri:"id" binding:"required"`
}

type UpdateRoleRequest struct {
	ID   uint   `json:"id" binding:"required" validate:"required"`
	Name string `json:"name" binding:"required" validate:"required"`
}

type DeleteRoleRequest struct {
	ID uint `json:"id" binding:"required" validate:"required"`
}

func GetRoles() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var roles []models.Role
		database.DB.Find(&roles)

		ctx.JSON(http.StatusOK, gin.H{
			"roles":   roles,
			"success": true,
		})
	}
}

func GetRole() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request GetRoleRequest
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.id_required",
				"success": false,
			})
			return
		}

		var role models.Role
		err := database.DB.Where("id = ?", request.ID).First(&role).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.role_not_found",
				"success": false,
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"role":    role,
			"success": true,
		})
	}
}

type CreateRoleRequest struct {
	Name string `json:"name" binding:"required"`
}

func CreateRole() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request CreateRoleRequest

		if err := ctx.BindJSON(&request); err != nil {
			log.Println(err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Error binding JSON: %v", err.Error()),
				"success": false,
			})
			return
		}

		role := models.Role{
			Name: request.Name,
		}

		dbc := database.DB.Create(&role)
		if dbc.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.unable_to_create_role",
				"success": false,
			})
			return
		}

			roleConfigs := []models.RoleConfig{
				{
					RoleId:    role.ID,
					Name:      "access_admin",
					Value:     "false",
					ValueType: "bool",
				},
				{
					RoleId:    role.ID,
					Name:      "access_home",
					Value:     "true",
					ValueType: "bool",
				},
				{
					RoleId:    role.ID,
					Name:      "header.show_notifications",
					Value:     "true",
					ValueType: "bool",
				},
				{
					RoleId:    role.ID,
					Name:      "header.show_profile",
					Value:     "true",
					ValueType: "bool",
				},
		}

		database.DB.Create(&roleConfigs)

		ctx.JSON(http.StatusOK, gin.H{
			"role":    role,
			"success": true,
		})
	}
}

func UpdateRole() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request UpdateRoleRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		var role models.Role
		err := database.DB.Where("id = ?", request.ID).First(&role).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.role_not_found",
				"success": false,
			})
		}

		err = database.DB.Model(&role).Update("name", request.Name).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"success": false,
			})
		}

		ctx.JSON(http.StatusOK, gin.H{
			"role":    role,
			"success": true,
		})
	}
}

func DeleteRole() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request DeleteRoleRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		var role models.Role
		err := database.DB.Where("id = ?", request.ID).First(&role).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.role_not_found",
				"success": false,
			})
		}

		err = database.DB.Delete(&role).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"success": false,
			})
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}
