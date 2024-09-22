package request

type (
	// AccountCreateRequest represents the payload required to create a new user account.
	// It includes the user's FullName, Email, and Password, all of which are mandatory fields.
	AccountCreateRequest struct {
		FullName string `json:"full_name" validate:"required,min=3,max=100"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=100"`
	}

	// AccountUpdateRequest represents a request to update account information such as FullName and Email.
	AccountUpdateRequest struct {
		FullName string `json:"full_name" validate:"required,min=3,max=100"`
		Email    string `json:"email" validate:"required,email"`
	}

	// AccountAuthRequest represents a request for authenticating an account.
	// It contains the necessary fields for Email and Password validation.
	AccountAuthRequest struct {
		Email    string `json:"email" validate:"required,email"`
		FcmToken string `json:"fcm_token" validate:"required"`
		Password string `json:"password" validate:"required,min=8,max=100"`
	}

	// AccountUnauthRequest represents the request payload for unauthenticating an account by invalidating access and refresh tokens.
	AccountUnauthRequest struct {
		AccessToken  string `json:"access_token" validate:"required"`
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	// AccountUpdatePasswordRequest represents a request to update an account password.
	AccountUpdatePasswordRequest struct {
		OldPassword string `json:"old_password" validate:"required,min=8,max=100"`
		NewPassword string `json:"new_password" validate:"required,min=8,max=100"`
	}

	// AccountRefreshTokenRequest represents a request to refresh an authentication token for an account.
	AccountRefreshTokenRequest struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
)
