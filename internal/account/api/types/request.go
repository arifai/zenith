package types

// CreateAccountRequest represents the payload required to create a new user account.
// It includes the user's FullName, Email, and Password, all of which are mandatory fields.
type CreateAccountRequest struct {
	FullName string `json:"full_name" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

// UpdateAccountRequest represents a request to update account information such as FullName and Email.
type UpdateAccountRequest struct {
	FullName string `json:"full_name,omitempty" validate:"omitempty,min=3,max=100"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
}

// AccountAuthRequest represents a request for authenticating an account.
// It contains the necessary fields for Email and Password validation.
type AccountAuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}
