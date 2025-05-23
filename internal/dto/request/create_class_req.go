package request

type CreateClassRequest struct {
	Code   string `json:"code" binding:"required,min=2,max=20"`
	Course string `json:"course" binding:"required,min=4,max=9"` // Ví dụ: 2021 hoặc "2021-2025"
}

var ClassValidateMessages = map[string]map[string]string{
	"Code": {
		"required": "Mã lớp là bắt buộc",
		"min":      "Mã lớp phải có ít nhất 2 ký tự",
		"max":      "Mã lớp không được vượt quá 20 ký tự",
	},
	"Course": {
		"required": "Khoá học là bắt buộc",
		"min":      "Khoá học phải có ít nhất 4 ký tự",
		"max":      "Khoá học không được vượt quá 9 ký tự",
	},
}
