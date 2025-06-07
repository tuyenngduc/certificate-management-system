package common

import "errors"

var (
	ErrUserNotExisted                 = errors.New("user_not_exists")
	ErrInvalidUserID                  = errors.New("invalid_user_id")
	ErrStudentIDExists                = errors.New("student_id_exists")
	ErrEmailExists                    = errors.New("email_exists")
	ErrUniversityNameExists           = errors.New("university_name_exists")
	ErrUniversityEmailDomainExists    = errors.New("university_email_domain_exists")
	ErrUniversityCodeExists           = errors.New("university_code_exists")
	ErrUniversityNotFound             = errors.New("university not found")
	ErrAccountUniversityNotFound      = errors.New("university account not found")
	ErrUniversityAlreadyApproved      = errors.New("university_already_approved")
	ErrAccountUniversityAlreadyExists = errors.New("university_admin_account_already_exists")
	ErrAccountNotFound                = errors.New("account_not_found")
	ErrInvalidOldPassword             = errors.New("invalid_old_password")
)
