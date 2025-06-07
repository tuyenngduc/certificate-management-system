package main

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("courseyear", func(fl validator.FieldLevel) bool {
			re := regexp.MustCompile(`^\d{4}$`)
			return re.MatchString(fl.Field().String())
		})

	}
}
