package request

type CreateSubjectRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Credit      int    `json:"credit" binding:"required,min=1"`
	FacultyID   string `json:"faculty_id" binding:"required"`
	Description string `json:"description"`
}

type UpdateSubjectRequest struct {
	Code        *string `json:"code"`
	Name        *string `json:"name"`
	Credit      *int    `json:"credit"`
	FacultyID   *string `json:"faculty_id"`
	Description *string `json:"description"`
}

type CreateSubjectByExcelRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Credit      int    `json:"credit" binding:"required,min=1"`
	FacultyCode string `json:"faculty_code" binding:"required"`
	Description string `json:"description"`
}

var SubjectValidateMessages = map[string]map[string]string{
	"Code": {
		"required": "Mã môn học là bắt buộc",
	},
	"Name": {
		"required": "Tên môn học là bắt buộc",
	},
	"Credit": {
		"required": "Số tín chỉ là bắt buộc",
		"min":      "Số tín chỉ phải lớn hơn 0",
	},
	"FacultyID": {
		"required": "ID khoa là bắt buộc",
	},
	"FacultyCode": {
		"required": "Mã khoa là bắt buộc",
	},
}
