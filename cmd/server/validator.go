package main

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/tuyenngduc/certificate-management-system/internal/models"
)

func InitValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("courseyear", func(fl validator.FieldLevel) bool {
			re := regexp.MustCompile(`^\d{4}$`)
			return re.MatchString(fl.Field().String())
		})

		v.RegisterStructValidation(models.ValidateCreateCertificateRequest, models.CreateCertificateRequest{})

		println("Đăng ký certtype validator")

	}
}
