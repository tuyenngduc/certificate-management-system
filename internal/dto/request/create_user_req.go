package request

type CreateUserRequest struct {
	FullName     string `json:"full_name" binding:"required,min=3,max=100"`
	Email        string `json:"email" binding:"omitempty,email"`
	StudentID    string `json:"student_id" binding:"required"`
	Ethnicity    string `json:"ethnicity" binding:"omitempty"`
	Gender       string `json:"gender" binding:"omitempty,oneof=Nam Nữ Khác"`
	FacultyCode  string `json:"faculty_code" binding:"required"`
	ClassCode    string `json:"class_code" binding:"required"`
	Course       string `json:"course" binding:"required"`
	NationalID   string `json:"national_id" binding:"omitempty,len=12,numeric"`
	Address      string `json:"address" binding:"omitempty"`
	PlaceOfBirth string `json:"place_of_birth" binding:"omitempty"`
	DateOfBirth  string `json:"date_of_birth" binding:"omitempty"` // dd/mm/yyyy
	PhoneNumber  string `json:"phone_number" binding:"omitempty,e164"`
}

var ValidateMessages = map[string]map[string]string{
	"FullName": {
		"required": "Họ tên là bắt buộc",
		"min":      "Họ tên phải có ít nhất 3 ký tự",
		"max":      "Họ tên không được vượt quá 100 ký tự",
	},
	"StudentID": {
		"required": "Mã sinh viên là bắt buộc",
	},
	"FacultyCode": {
		"required": "Mã khoa là bắt buộc",
	},
	"ClassCode": {
		"required": "Mã lớp là bắt buộc",
	},
	"Course": {
		"required": "Khoá học là bắt buộc",
	},
	"NationalID": {
		"len":     "CCCD/CMND phải đúng 12 số",
		"numeric": "CCCD/CMND chỉ được chứa ký tự số",
	},
	"Email": {
		"email": "Email không hợp lệ",
	},
	"Gender": {
		"oneof": "Giới tính phải là Nam, Nữ hoặc Khác",
	},
	"PhoneNumber": {
		"e164": "Số điện thoại không đúng định dạng quốc tế",
	},
}
