package response

type ClassResponse struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Course      string `json:"course"`
	FacultyCode string `json:"faculty_code"`
	FacultyName string `json:"faculty_name"`
}
