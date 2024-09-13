package types

// CreateAccountRequest is a struct that represent the request body for creating an account
type CreateAccountRequest struct {
	FullName string `json:"full_name" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

// UpdateAccountRequest is a struct that represent the request body for updating an account
type UpdateAccountRequest struct {
	FullName string `json:"full_name,omitempty" validate:"omitempty,min=3,max=100"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
}

// AccountAuthRequest is a struct that represent the request body for authorizing an account
type AccountAuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}
