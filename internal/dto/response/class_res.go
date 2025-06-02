package response

type ClassResponse struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Course      string `json:"course"`
	FacultyID   string `json:"faculty_id"`
	FacultyCode string `json:"faculty_code"`
	FacultyName string `json:"faculty_name"`
}
