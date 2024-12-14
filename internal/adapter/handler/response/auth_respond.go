package response

type SuccessAuthResponse struct {
	Meta
	AccessToken string `json:"access_token"`
	ExpiresAt string `json:"expires_at"`
}