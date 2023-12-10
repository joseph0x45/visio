package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/sirupsen/logrus"
)

func ExchangeCodeWithToken(code string, logger *logrus.Logger) (string, error) {
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
		logger.Error(err)
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	respBodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Error(err)
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
		logger.Error(err)
		return "", err
	}
	return githubResponse.AccessToken, nil
}

func GetUserGithubData(access_token string, logger *logrus.Logger) (map[string]interface{}, error) {
	req, err := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", access_token))
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	respBodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	data := map[string]interface{}{}
	err = json.Unmarshal(
		respBodyBytes,
		&data,
	)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	req, err = http.NewRequest(
		"GET",
		"https://api.github.com/user/emails",
		nil,
	)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", access_token))
	req.Header.Set("Accept", "application/json")
	response, err = http.DefaultClient.Do(req)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	respBodyBytes, err = io.ReadAll(response.Body)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
  println(string(respBodyBytes))
	return data, nil
}

func GetGithubUserEmail(token string, logger *logrus.Logger) (email string, err error) {
	req, err := http.NewRequest(
		"GET",
		"https://api.github.com/user/emails",
		nil,
	)
	if err != nil {
		logger.Error(err)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	req.Header.Set("Accept", "application/json")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error(err)
		return
	}
	respBytes, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Error(err)
		return
	}
	emails := []struct {
		Email   string `json:"email"`
		Primary bool   `json:"primary"`
	}{}
	err = json.Unmarshal(
		respBytes,
		&emails,
	)
	if err != nil {
		logger.Error(err)
		return
	}
	for _, userEmail := range emails {
		if userEmail.Primary == true {
			email = userEmail.Email
			break
		}
	}
	return
}
