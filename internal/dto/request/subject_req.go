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
