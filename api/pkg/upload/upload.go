package upload

import (
	"fmt"
	"net/http"
	"time"
	"path/filepath"

	"github.com/berkaycubuk/billiard_software_api/database"
	"github.com/berkaycubuk/billiard_software_api/middlewares"
	"github.com/berkaycubuk/billiard_software_api/models"
	"github.com/berkaycubuk/billiard_software_api/utils"
	"github.com/gin-gonic/gin"
)

func uploadFile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		file, err := ctx.FormFile("file")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"success": false,
			})
			return
		}

		extension := filepath.Ext(file.Filename)

		filename := file.Filename 
		filenameFormatted := fmt.Sprintf("%v-%v", (time.Now().Unix()), utils.RandomString(5)) + extension
		fmt.Println(filenameFormatted)

		destination := "static/uploads/" + filename

		uploadModel := models.Upload{
			Size:     uint(file.Size),
			MimeType: "image/png",
			FileName: filename,
		}

		if err := ctx.SaveUploadedFile(file, destination); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"success": false,
			})
			return
		}

		err = database.DB.Create(&uploadModel).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"success": false,
			})
		}

		ctx.JSON(http.StatusOK, uploadModel.ID)
	}
}

func Routes(r *gin.RouterGroup) {
	r.Use(middlewares.AuthMiddleware())
	r.POST("/upload", uploadFile())
}
