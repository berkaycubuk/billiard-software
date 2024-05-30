package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/berkaycubuk/billiard_software_api/database"
	"github.com/berkaycubuk/billiard_software_api/models"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func GetProducts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var products []models.Product
		database.DB.Where("deleted_at IS NULL").
			Order("`order` asc").
			Preload("Image").Find(&products)

		ctx.JSON(http.StatusOK, gin.H{
			"products": products,
			"success":  true,
		})
	}
}

type GetProductRequest struct {
	ID string `uri:"id" binding:"required"`
}

func GetProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request GetProductRequest
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "ID required.",
				"success": false,
			})
			return
		}

		var product models.Product
		dbs := database.DB.Where("id = ?", request.ID).Preload("Image").First(&product)
		if dbs.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.unable_to_find_product",
				"success": false,
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"product": product,
			"success": true,
		})
	}
}

type UpdateProductRequest struct {
	ID    uint            `json:"id" binding:"required" validate:"required"`
	Name  string          `json:"name" binding:"required" validate:"required"`
	Price decimal.Decimal `json:"price" binding:"required" validate:"required"`
	Image uint            `json:"image"`
}

func UpdateProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request UpdateProductRequest

		if err := ctx.BindJSON(&request); err != nil {
			log.Println(err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Error binding JSON: %v", err.Error()),
				"success": false,
			})
			return
		}

		var product models.Product
		dbs := database.DB.Where("id = ?", request.ID).First(&product)
		if dbs.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.unable_to_find_product",
				"success": false,
			})
			return
		}

		database.DB.Model(&product).Updates(request)

		if request.Image == uint(0) {
			ctx.JSON(http.StatusOK, gin.H{
				"product": product,
				"success": true,
			})
			return
		}

		// check is image upload successfull
		var upload models.Upload
		err := database.DB.Where("id = ?", request.Image).First(&upload).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.image_upload_not_successful",
				"success": false,
			})
			return
		}

		// check are there any image before
		var productImage models.ProductImage
		err = database.DB.Where("product_id = ?", request.ID).First(&productImage).Error
		if err == nil {
			// update current image
			database.DB.Model(&productImage).Updates(models.ProductImage{
				UploadID:       request.Image,
				UploadFilename: upload.FileName,
				UpdatedAt:      time.Now(),
			})
		} else {
			// create a new product image
			database.DB.Create(&models.ProductImage{
				ProductID:      product.ID,
				UploadID:       request.Image,
				UploadFilename: upload.FileName,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			})
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

type UpdateProductOrderRequest struct {
	ID uint `json:"id" binding:"required" validate:"required"`
	Direction int `json:"direction" binding:"required" validate:"required"`
}
// 1 -> up
// 2 -> down

func UpdateProductOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request UpdateProductOrderRequest

		if err := ctx.BindJSON(&request); err != nil {
			log.Println(err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Error binding JSON: %v", err.Error()),
				"success": false,
			})
			return
		}

		var product models.Product
		err := database.DB.Where("id = ?", request.ID).First(&product).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.unable_to_find_product",
				"success": false,
			})
			return
		}

		if request.Direction == 1 {
			if product.Order == 1 {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": "errors.product_at_top",
					"success": false,
				})
				return
			}

			newOrderNumber := product.Order - 1

			var tempProduct models.Product
			err := database.DB.Where("`order` = ?", newOrderNumber).First(&tempProduct).Error
			if err == nil {
				log.Println(tempProduct.ID)
				database.DB.Model(&tempProduct).Update("order", newOrderNumber + 1)
			}
			log.Println(err)

			database.DB.Model(&product).Update("order", newOrderNumber)

			ctx.JSON(http.StatusOK, gin.H{
				"success": true,
			})
			return
		} else {
			newOrderNumber := product.Order + 1

			var tempProduct models.Product
			err := database.DB.Where("`order` = ?", newOrderNumber).First(&tempProduct).Error
			if err == nil {
				log.Println(tempProduct.ID)
				database.DB.Model(&tempProduct).Update("order", newOrderNumber - 1)
			}
			log.Println(err)

			database.DB.Model(&product).Update("order", newOrderNumber)

			ctx.JSON(http.StatusOK, gin.H{
				"success": true,
			})
			return
		}
	}
}

type DeleteProductRequest struct {
	ID uint `json:"id" binding:"required"`
}

func DeleteProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request DeleteProductRequest

		if err := ctx.BindJSON(&request); err != nil {
			log.Println(err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Error binding JSON: %v", err.Error()),
				"success": false,
			})
			return
		}

		var product models.Product
		dbs := database.DB.Where("id = ?", request.ID).First(&product)
		if dbs.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.unable_to_find_product",
				"success": false,
			})
			return
		}

		// instead of just deleting it, set deleted_at
		err := database.DB.Model(&product).Update("deleted_at", time.Now()).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"success": false,
			})
			return
		}
		//database.DB.Delete(&product)

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

type CreateProductRequest struct {
	Name  string          `json:"name" binding:"required" validate:"required"`
	Price decimal.Decimal `json:"price" binding:"required" validate:"required"`
	Image uint            `json:"image"`
}

func CreateProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request CreateProductRequest

		if err := ctx.BindJSON(&request); err != nil {
			log.Println(err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Error binding JSON: %v", err.Error()),
				"success": false,
			})
			return
		}

		var orderValue int

		// get last order value
		var products []models.Product
		err := database.DB.Find(&products).Error
		if err != nil {
			orderValue = 1
		} else {
			for _, v := range products { // find highest order
				if v.Order > orderValue {
					orderValue = v.Order
				}
			}

			orderValue = orderValue + 1
		}

		// create product model
		product := models.Product{
			Name:   request.Name,
			Price:  request.Price,
			Order:  orderValue,
			Status: models.ON_SALE,
		}

		// create product row
		dbc := database.DB.Create(&product)
		if dbc.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.unable_to_create_product",
				"success": false,
			})
			return
		}

		if request.Image == uint(0) {
			ctx.JSON(http.StatusOK, gin.H{
				"product": product,
				"success": true,
			})
			return
		}

		// check is image upload successfull
		var upload models.Upload
		err = database.DB.Where("id = ?", request.Image).First(&upload).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.image_upload_not_successfull",
				"success": false,
			})
			return
		}

		// create a new product image
		database.DB.Create(&models.ProductImage{
			ProductID:      product.ID,
			UploadID:       request.Image,
			UploadFilename: upload.FileName,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		})

		ctx.JSON(http.StatusOK, gin.H{
			"product": product,
			"success": true,
		})
	}
}
