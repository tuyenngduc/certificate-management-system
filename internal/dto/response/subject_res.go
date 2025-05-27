package response

type SubjectResponse struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Credit      int    `json:"credit"`
	FacultyID   string `json:"faculty_name"`
	Description string `json:"description,omitempty"`
}
