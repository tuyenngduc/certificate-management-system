package common

import "errors"

var (
	ErrStudentIDExists             = errors.New("student_id_exists")
	ErrEmailExists                 = errors.New("email_exists")
	ErrUserNotExisted              = errors.New("user_not_exists")
	ErrInvalidUserID               = errors.New("invalid_user_id")
	ErrUniversityNameExists        = errors.New("university_name_exists")
	ErrUniversityEmailDomainExists = errors.New("university_email_domain_exists")
	ErrUniversityCodeExists        = errors.New("university_code_exists")
	ErrUniversityNotFound          = errors.New("university not found")
)
