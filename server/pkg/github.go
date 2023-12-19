package pkg

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
  "fmt"
)

func ExchangeCodeWithToken(code string) (string, error) {
	data := url.Values{}
	data.Set("client_id", os.Getenv("GITHUB_CLIENT_ID"))
	data.Set("client_secret", os.Getenv("GITHUB_CLIENT_SECRET"))
	data.Set("code", code)
	req, err := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBufferString(data.Encode()),
	)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	respBodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	githubResponse := struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}{}
	err = json.Unmarshal(
		respBodyBytes,
		&githubResponse,
	)
	if err != nil {
		return "", err
	}
	return githubResponse.AccessToken, nil
}

func GetUserGithubData(access_token string) (map[string]interface{}, error) {
	req, err := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", access_token))
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	respBodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	data := map[string]interface{}{}
	err = json.Unmarshal(
		respBodyBytes,
		&data,
	)
	if err != nil {
		return nil, err
	}
	return data, nil
}
