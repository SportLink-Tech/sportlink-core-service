package request

type GoogleAuthRequest struct {
	IDToken string `json:"id_token" validate:"required"`
}
