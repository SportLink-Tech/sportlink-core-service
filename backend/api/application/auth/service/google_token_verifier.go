package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type GoogleTokenInfo struct {
	Sub              string `json:"sub"`
	Email            string `json:"email"`
	EmailVerified    string `json:"email_verified"`
	Name             string `json:"name"`
	Picture          string `json:"picture"`
	Aud              string `json:"aud"`
	ErrorDescription string `json:"error_description,omitempty"`
}

type GoogleTokenVerifier interface {
	Verify(ctx context.Context, idToken string) (*GoogleTokenInfo, error)
}

type httpGoogleTokenVerifier struct {
	clientID string
}

func NewGoogleTokenVerifier(clientID string) GoogleTokenVerifier {
	return &httpGoogleTokenVerifier{clientID: clientID}
}

func (v *httpGoogleTokenVerifier) Verify(_ context.Context, idToken string) (*GoogleTokenInfo, error) {
	resp, err := http.Get("https://oauth2.googleapis.com/tokeninfo?id_token=" + idToken)
	if err != nil {
		return nil, fmt.Errorf("failed to contact Google: %w", err)
	}
	defer resp.Body.Close()

	var info GoogleTokenInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("failed to decode Google response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid token: %s", info.ErrorDescription)
	}

	if info.EmailVerified != "true" {
		return nil, fmt.Errorf("email not verified by Google")
	}

	if info.Aud != v.clientID {
		return nil, fmt.Errorf("token audience mismatch")
	}

	return &info, nil
}
