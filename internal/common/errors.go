package common

import "errors"

var (
	ErrStudentIDExists = errors.New("student_id_exists")
	ErrEmailExists     = errors.New("email_exists")
)
