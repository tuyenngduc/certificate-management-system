package request

type BulkCreateUserRequest struct {
	Users []CreateUserRequest `json:"users" binding:"required,dive"`
}
