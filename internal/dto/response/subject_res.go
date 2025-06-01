package response

type SubjectResponse struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Credit      int    `json:"credit"`
	FacultyCode string `json:"faculty_code"`
	FacultyName string `json:"faculty_name"`
	Description string `json:"description,omitempty"`
}
