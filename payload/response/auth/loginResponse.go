package auth

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refereshToken"`
}
