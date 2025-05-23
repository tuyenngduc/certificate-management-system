package request

type CreateFacultyRequest struct {
	Name           string `json:"name" binding:"required,min=3,max=100"`
	Code           string `json:"code" binding:"required,min=2,max=20"`
	TrainingPeriod string `json:"training_period" binding:"required,oneof='3.5 năm' '4 năm' '4.5 năm' '5 năm'"`
}

var FacultyValidateMessages = map[string]map[string]string{
	"Name": {
		"required": "Tên khoa là bắt buộc",
		"min":      "Tên khoa phải có ít nhất 3 ký tự",
		"max":      "Tên khoa không được vượt quá 100 ký tự",
	},
	"Code": {
		"required": "Mã khoa là bắt buộc",
		"min":      "Mã khoa phải có ít nhất 2 ký tự",
		"max":      "Mã khoa không được vượt quá 20 ký tự",
	},
	"TrainingPeriod": {
		"required": "Thời gian đào tạo là bắt buộc",
		"oneof":    "Thời gian đào tạo chỉ được chọn: 3.5 năm, 4 năm, 4.5 năm, 5 năm",
	},
}
