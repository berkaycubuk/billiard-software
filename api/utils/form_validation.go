package utils

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
    snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
    snake  = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
    return strings.ToLower(snake)
}

type ErrorMsg struct {
	Field	string	`json:"field"`
	Message	string	`json:"message"`
}

func GetErrorMessage(e validator.FieldError) string {
	var message string

	switch e.Tag() {
	case "required":
		message = "required"
		break
	}

	return message
}

func ValidateRequest(c *gin.Context, request any) bool {
	if err := c.ShouldBindJSON(request); err != nil {
		v := validator.New()

		err := v.Struct(request)
		errors := make(map[string]string, len(err.(validator.ValidationErrors)))
		for _, e := range err.(validator.ValidationErrors) {
			errors[ToSnakeCase(e.Field())] = GetErrorMessage(e)
		}

		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": errors,
			"success": false,
		})
		return false
	}

	return true
}
