package lib

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Scotiacon-Tech/libs/message-relay/go/requests"
	"github.com/gofiber/fiber/v2"
)

var KeyInvalidError = errors.New("Key Invalid")

func (client Client) RequestSend(key string, service string, req *requests.SendRequest) (*requests.SendResponse, error) {
	if key == "" {
		return nil, KeyInvalidError
	}

	url := fmt.Sprintf("%s/send/%s", client.Config.ServerURL, service)

	reqBody, _ := json.Marshal(req)

	agent := fiber.Post(url)
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
	req := requests.TokenRequest{
		GrantType:    "client_credentials",
		ClientID:     client.Config.ClientID,
		ClientSecret: client.Config.ClientSecret,
		Audience:     []string{client.Config.AudienceUUID},
		Scope:        "openid",
	}

	reqBody, _ := json.Marshal(req)

	agent := fiber.Post(client.Config.TokenEndpoint)
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
	url := fmt.Sprintf("%s/auth", client.Config.ServerURL)

	agent := fiber.Post(url)
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
