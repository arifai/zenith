package response

type (
	// AccountAuthResponse represents the structure of the authentication response containing AccessToken and RefreshToken.
	AccountAuthResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
)
