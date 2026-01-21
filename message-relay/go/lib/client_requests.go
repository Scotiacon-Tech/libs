package lib

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/Scotiacon-Tech/libs/message-relay/go/requests"
	"github.com/gofiber/fiber/v2"
)

var KeyInvalidError = errors.New("Key Invalid")

func (client Client) RequestSend(key string, service string, to any, subject string, body string) (*requests.SendResponse, error) {
	if key == "" {
		return nil, KeyInvalidError
	}

	url := os.Getenv("SERVER_URL")

	req := requests.SendRequest{
		To:      to,
		Subject: subject,
		Body:    body,
	}

	reqBody, _ := json.Marshal(req)

	agent := fiber.Post(url + "/send/" + service)
	agent.Request().Header.Set("Authorization", "Bearer "+key)
	agent.Request().Header.Set("Content-Type", "application/json")
	agent.Body(reqBody)

	code, res, errs := agent.Bytes()

	if code == 401 {
		return nil, KeyInvalidError
	} else if len(errs) > 0 || code != 200 {
		return nil, errors.New("Request failed")
	}

	var sendRes requests.SendResponse
	err := json.Unmarshal(res, &sendRes)

	if err != nil {
		return nil, errors.New("Failed to decode JSON")
	}

	return &sendRes, nil
}

func (client Client) RequestJWT() (*requests.TokenResponse, error) {
	url := os.Getenv("TOKEN_ENDPOINT")

	req := requests.TokenRequest{
		GrantType:    "client_credentials",
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Audience:     []string{os.Getenv("AUDIENCE_UUID")},
		Scope:        "openid",
	}

	reqBody, _ := json.Marshal(req)

	agent := fiber.Post(url)
	agent.Request().Header.Set("Content-Type", "application/json")
	agent.Body(reqBody)

	code, res, errs := agent.Bytes()

	if len(errs) > 0 || code != 200 {
		return nil, errors.New("Request failed")
	}

	var tokenRes requests.TokenResponse
	err := json.Unmarshal(res, &tokenRes)

	if err != nil {
		return nil, errors.New("Failed to decode JSON")
	}

	return &tokenRes, nil
}

func (client Client) RequestKey(jwt string) (*requests.KeyResponse, error) {
	url := os.Getenv("SERVER_URL")

	agent := fiber.Post(url + "/auth")
	agent.Request().Header.Set("Authorization", "Bearer "+jwt)

	code, res, errs := agent.Bytes()

	if len(errs) > 0 || code != 200 {
		return nil, errors.New("Request failed")
	}

	var keyRes requests.KeyResponse
	err := json.Unmarshal(res, &keyRes)

	if err != nil {
		return nil, errors.New("Failed to decode JSON")
	}

	return &keyRes, nil
}
