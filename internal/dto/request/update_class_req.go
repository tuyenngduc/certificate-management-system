package request

type UpdateClassRequest struct {
	Code   string `json:"code" binding:"omitempty,min=2,max=20"`
	Course string `json:"course" binding:"omitempty,min=4,max=9"`
}

var ClassUpdateValidateMessages = map[string]map[string]string{
	"Code": {
		"min": "Mã lớp phải có ít nhất 2 ký tự",
		"max": "Mã lớp không được vượt quá 20 ký tự",
	},
	"Course": {
		"min": "Khoá học phải có ít nhất 4 ký tự",
		"max": "Khoá học không được vượt quá 9 ký tự",
	},
}
