package request

type CreateLecturerRequest struct {
	Code      string `json:"code" binding:"required,min=3,max=20"`
	FullName  string `json:"full_name" binding:"required,min=3,max=100"`
	Email     string `json:"email" binding:"required,email"`
	Title     string `json:"title" binding:"required,oneof=ThS TS PGS GS"`
	FacultyID string `json:"faculty_id" binding:"required"` // ID khoa (ObjectID dạng string)
}

var LecturerValidateMessages = map[string]map[string]string{
	"Code": {
		"required": "Mã giảng viên là bắt buộc",
		"min":      "Mã giảng viên phải có ít nhất 3 ký tự",
		"max":      "Mã giảng viên không được vượt quá 20 ký tự",
	},
	"FullName": {
		"required": "Họ tên là bắt buộc",
		"min":      "Họ tên phải có ít nhất 3 ký tự",
		"max":      "Họ tên không được vượt quá 100 ký tự",
	},
	"Email": {
		"required": "Email là bắt buộc",
		"email":    "Email không hợp lệ",
	},
	"Title": {
		"required": "Chức danh là bắt buộc",
		"oneof":    "Chức danh phải là một trong: ThS, TS, PGS, GS",
	},
	"FacultyID": {
		"required": "ID khoa là bắt buộc",
	},
}
