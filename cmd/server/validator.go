package main

import (
	"regexp"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InitValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("objectid", func(fl validator.FieldLevel) bool {
			_, err := primitive.ObjectIDFromHex(fl.Field().String())
			return err == nil
		})

		v.RegisterValidation("not_future", func(fl validator.FieldLevel) bool {
			t, ok := fl.Field().Interface().(time.Time)
			return ok && !t.After(time.Now())
		})
		v.RegisterValidation("courseyear", func(fl validator.FieldLevel) bool {
			re := regexp.MustCompile(`^\d{4}-\d{4}$`)
			return re.MatchString(fl.Field().String())
		})

	}
}
