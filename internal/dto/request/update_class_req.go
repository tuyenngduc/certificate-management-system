package request

type UpdateClassRequest struct {
	Code      string `json:"code" binding:"omitempty,min=2,max=20"`
	Course    string `json:"course" binding:"omitempty,courseyear"`
	FacultyID string `json:"faculty_id" binding:"omitempty"`
}

var ClassUpdateValidateMessages = map[string]map[string]string{
	"Code": {
		"min": "Mã lớp phải có ít nhất 2 ký tự",
		"max": "Mã lớp không được vượt quá 20 ký tự",
	},
	"Course": {
		"courseyear": "Khoá học phải có định dạng năm-năm, ví dụ: 2021-2025",
	},
}
