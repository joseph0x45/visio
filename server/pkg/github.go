package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"visio/internal/types"
)

type UserData struct {
	Avatar string  `json:"avatar_url"`
	Login  string  `json:"login"`
	Id     float64 `json:"id"`
}

func GetToken(code string) (string, error) {
	formData := url.Values{}
	formData.Set("client_id", os.Getenv("GH_CLIENT_ID"))
	formData.Set("client_secret", os.Getenv("GH_CLIENT_SECRET"))
	formData.Set("code", code)
	req, err := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBufferString(formData.Encode()),
	)
	if err != nil {
		return "", fmt.Errorf("Error while creating http request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error while sending http request: %w", err)
	}
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("Error while reading request body: %w", err)
	}
	githubResponse := struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}{}
	err = json.Unmarshal(responseBody, &githubResponse)
	if err != nil {
		return "", fmt.Errorf("Error while unmarshalling body: %w", err)
	}
	return githubResponse.AccessToken, nil
}

func GetUserData(token string) (*UserData, error) {
  userData := new(UserData)
	req, err := http.NewRequest(
		"GET",
		"http://api.github.com/user",
		nil,
	)
	if err != nil {
		return userData, fmt.Errorf("Error while creating http request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return userData, fmt.Errorf("Error while sending http request: %w", err)
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return userData, fmt.Errorf("Error while reading response body: %w", err)
	}
	err = json.Unmarshal(responseBody, userData)
	if err != nil {
		return nil, fmt.Errorf("Error while unmarshalling body: %w", err)
	}
	return userData, nil
}

func GetUserPrimaryEmail(token string) (string, error) {
	req, err := http.NewRequest(
		"GET",
		"https://api.github.com/user/emails",
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("Error while creating http request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/json")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error while sending http request: %w", err)
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("Error while reading response body: %w", err)
	}
	emailData := []struct {
		Email      string `json:"email"`
		Primary    bool   `json:"primary"`
		Verified   bool   `json:"verified"`
		Visibility string `json:"visibility"`
	}{}
	err = json.Unmarshal(responseBody, &emailData)
	if err != nil {
		return "", fmt.Errorf("Error while unmarshalling body: %w", err)
	}
	primaryEmail := ""
	for _, email := range emailData {
		if email.Primary {
			primaryEmail = email.Email
			break
		}
	}
	if primaryEmail == "" {
		return "", types.ErrNoPrimaryEmailFound
	}
	return primaryEmail, nil
}
