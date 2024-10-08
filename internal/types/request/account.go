package request

type (
	// AccountCreateRequest represents the payload required to create a new user account.
	// It includes the user's FullName, Email, and Password, all of which are mandatory fields.
	AccountCreateRequest struct {
		FullName string `json:"full_name" validate:"required,min=3,max=100" reason:"required:Full name is required"`
		Email    string `json:"email" validate:"required,email" reason:"required:Email is required;email:Invalid email address"`
		Password string `json:"password" validate:"required,min=8,max=100" reason:"required:Password is required;min:Password must be at least 8 characters;max:Password must be at most 100 characters"`
	}

	// AccountUpdateRequest represents a request to update account information such as FullName and Email.
	AccountUpdateRequest struct {
		FullName string `json:"full_name" validate:"required,min=3,max=100" reason:"required:Full name is required;min:Full name must be at least 3 characters;max:Full name must be at most 100 characters"`
		Email    string `json:"email" validate:"required,email" reason:"required:Email is required;email:Invalid email address"`
	}

	// AccountAuthRequest represents a request for authenticating an account.
	// It contains the necessary fields for Email and Password validation.
	AccountAuthRequest struct {
		Email    string `json:"email" validate:"required,email" reason:"required:Email is required;email:Invalid email address"`
		FcmToken string `json:"fcm_token" validate:"required" reason:"required:FCM token is required"`
		Password string `json:"password" validate:"required,min=8,max=100" reason:"required:Password is required;min:Password must be at least 8 characters;max:Password must be at most 100 characters"`
		DeviceID string `json:"device_id" validate:"required,uuid" reason:"required:Device ID is required;uuid:Device ID must be a valid UUID"`
	}

	// AccountUnauthRequest represents the request payload for unauthenticating an account by invalidating access and refresh tokens.
	AccountUnauthRequest struct {
		AccessToken  string `json:"access_token" validate:"required" reason:"required:Access token is required"`
		RefreshToken string `json:"refresh_token" validate:"required" reason:"required:Refresh token is required"`
	}

	// AccountUpdatePasswordRequest represents a request to update an account password.
	AccountUpdatePasswordRequest struct {
		OldPassword string `json:"old_password" validate:"required,min=8,max=100" reason:"required:Old password is required;min:Old password must be at least 8 characters;max:Old password must be at most 100 characters"`
		NewPassword string `json:"new_password" validate:"required,min=8,max=100" reason:"required:New password is required;min:New password must be at least 8 characters;max:New password must be at most 100 characters"`
	}

	// AccountRefreshTokenRequest represents a request to refresh an authentication token for an account.
	AccountRefreshTokenRequest struct {
		RefreshToken string `json:"refresh_token" validate:"required" reason:"required:Refresh token is required"`
	}
)
