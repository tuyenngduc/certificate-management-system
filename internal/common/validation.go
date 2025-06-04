package common

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

func TranslateError(field, tag string) string {
	messages := map[string]map[string]string{
		"StudentID": {
			"required": "Mã sinh viên không được để trống",
		},
		"FullName": {
			"required": "Họ tên không được để trống",
		},
		"Email": {
			"required": "Email không được để trống",
			"email":    "Email không hợp lệ",
		},
		"Faculty": {
			"required": "Khoa không được để trống",
		},
		"Class": {
			"required": "Lớp không được để trống",
		},
		"Course": {
			"required":   "Khóa học không được để trống",
			"courseyear": "Khóa học phải có định dạng NĂM-NĂM, ví dụ 2022-2026",
		},
		"PersonalEmail": {
			"required": "Email cá nhân không được để trống",
			"email":    "Email cá nhân không hợp lệ",
		},
		"Password": {
			"required": "Mật khẩu không được để trống",
		},
	}

	if fieldMsg, ok := messages[field]; ok {
		if msg, ok2 := fieldMsg[tag]; ok2 {
			return msg
		}
	}
	return field + " không hợp lệ"
}

func ParseValidationError(err error) (map[string]string, bool) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		errs := make(map[string]string)
		for _, e := range ve {
			field := e.Field()
			tag := e.Tag()
			errs[field] = TranslateError(field, tag)
		}
		return errs, true
	}
	return nil, false
}

func ParseMongoDuplicateError(err error) string {
	if mongoErr, ok := err.(mongo.WriteException); ok {
		for _, writeErr := range mongoErr.WriteErrors {
			if writeErr.Code == 11000 {
				if strings.Contains(writeErr.Message, "studentId") {
					return "Mã sinh viên đã tồn tại"
				} else if strings.Contains(writeErr.Message, "email") {
					return "Email đã tồn tại"
				}
			}
		}
	}
	return ""
}
