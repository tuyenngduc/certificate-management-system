package request

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

type CreateClassRequest struct {
	Code      string `json:"code" binding:"required,min=2,max=20"`
	Course    string `json:"course" binding:"required,courseyear"`
	FacultyID string `json:"faculty_id" binding:"required"`
}

var ClassValidateMessages = map[string]map[string]string{
	"Code": {
		"required": "Mã lớp là bắt buộc",
		"min":      "Mã lớp phải có ít nhất 2 ký tự",
		"max":      "Mã lớp không được vượt quá 20 ký tự",
	},
	"Course": {
		"required":   "Khoá học là bắt buộc",
		"courseyear": "Khoá học phải có định dạng năm-năm, ví dụ: 2021-2025",
	},
	"FacultyID": {
		"required": "ID khoa là bắt buộc",
	},
}

func RegisterClassValidators(v *validator.Validate) {
	v.RegisterValidation("courseyear", func(fl validator.FieldLevel) bool {
		re := regexp.MustCompile(`^\d{4}-\d{4}$`)
		return re.MatchString(fl.Field().String())
	})
}
