package types

// AccountCreateRequest represents the payload required to create a new user account.
// It includes the user's FullName, Email, and Password, all of which are mandatory fields.
type AccountCreateRequest struct {
	FullName string `json:"full_name" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

// AccountUpdateRequest represents a request to update account information such as FullName and Email.
type AccountUpdateRequest struct {
	FullName string `json:"full_name,omitempty" validate:"omitempty,min=3,max=100"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
}

// AccountAuthRequest represents a request for authenticating an account.
// It contains the necessary fields for Email and Password validation.
type AccountAuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

// AccountUpdatePasswordRequest represents a request to update an account password.
type AccountUpdatePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required,min=8,max=100"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=100"`
}
