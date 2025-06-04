package common

import "errors"

var (
	ErrStudentIDExists = errors.New("student_id_exists")
	ErrEmailExists     = errors.New("email_exists")
	ErrUserNotExisted  = errors.New("user_not_exists")
	ErrInvalidUserID   = errors.New("invalid_user_id")
)
