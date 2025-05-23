package request

type UpdateLecturerRequest struct {
	Code     string `json:"code" binding:"omitempty,min=3,max=20"`
	FullName string `json:"full_name" binding:"omitempty,min=3,max=100"`
	Email    string `json:"email" binding:"omitempty,email"`
	Title    string `json:"title" binding:"omitempty,oneof=ThS TS PGS GS"`
}

var LecturerUpdateValidateMessages = map[string]map[string]string{
	"Code": {
		"min": "Mã giảng viên phải có ít nhất 3 ký tự",
		"max": "Mã giảng viên không được vượt quá 20 ký tự",
	},
	"Email": {
		"email": "Email không hợp lệ",
	},
	"Title": {
		"oneof": "Chức danh phải là một trong: ThS, TS, PGS, GS",
	},
}
