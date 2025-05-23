package response

type UserResponse struct {
	ID           string `json:"id"`
	StudentID    string `json:"student_id"`
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	Ethnicity    string `json:"ethnicity"`
	Gender       string `json:"gender"`
	FacultyCode  string `json:"faculty_code"`
	ClassCode    string `json:"class_code"`
	Course       string `json:"course"`
	NationalID   string `json:"national_id"`
	Address      string `json:"address"`
	PlaceOfBirth string `json:"place_of_birth"`
	DateOfBirth  string `json:"date_of_birth"`
	PhoneNumber  string `json:"phone_number,omitempty"`
}
