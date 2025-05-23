package request

type UpdateFacultyRequest struct {
	Name           string `json:"name" binding:"omitempty,min=3,max=100"`
	Code           string `json:"code" binding:"omitempty,min=2,max=20"`
	TrainingPeriod string `json:"training_period" binding:"omitempty,oneof='3.5 năm' '4 năm' '4.5 năm' '5 năm'"`
}

var FacultyUpdateValidateMessages = map[string]map[string]string{
	"Name": {
		"min": "Tên khoa phải có ít nhất 3 ký tự",
		"max": "Tên khoa không được vượt quá 100 ký tự",
	},
	"Code": {
		"min": "Mã khoa phải có ít nhất 2 ký tự",
		"max": "Mã khoa không được vượt quá 20 ký tự",
	},
	"TrainingPeriod": {
		"oneof": "Thời gian đào tạo chỉ được chọn: 3.5 năm, 4 năm, 4.5 năm, 5 năm",
	},
}
