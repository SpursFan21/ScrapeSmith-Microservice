package models

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IsAdmin      bool   `json:"is_admin"`
}
