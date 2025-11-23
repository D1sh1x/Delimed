package cdek

import (
	"bytes"
	"context"
	"delimed/internal/transport/dto/response"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const baseURL = "https://api.edu.cdek.ru/v2"

func GetCDEKToken(ctx context.Context, clientID, clientSecret string) (string, error) {
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseURL+"/oauth/token",
		bytes.NewBufferString(form.Encode()),
	)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("auth failed: %s", resp.Status)
	}

	var ar response.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&ar); err != nil {
		return "", err
	}

	if ar.AccessToken == "" {
		return "", fmt.Errorf("empty access_token in response")
	}

	return ar.AccessToken, nil
}
