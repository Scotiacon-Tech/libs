package lib

import (
	"encoding/json"
	"fmt"

	"github.com/Scotiacon-Tech/libs/message-relay/go/requests"
	"github.com/gofiber/fiber/v2"
)

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
		return nil, RequestFailedError
	}

	var sendRes requests.SendResponse
	err := json.Unmarshal(res, &sendRes)

	if err != nil {
		return nil, JSONDecodeError
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
		return nil, RequestFailedError
	}

	var tokenRes requests.TokenResponse
	err := json.Unmarshal(res, &tokenRes)

	if err != nil {
		return nil, JSONDecodeError
	}

	return &tokenRes, nil
}

func (client Client) RequestKey(jwt string) (*requests.KeyResponse, error) {
	url := fmt.Sprintf("%s/auth", client.Config.ServerURL)

	agent := fiber.Post(url)
	agent.Request().Header.Set("Authorization", "Bearer "+jwt)

	code, res, errs := agent.Bytes()

	if len(errs) > 0 || code != 200 {
		return nil, RequestFailedError
	}

	var keyRes requests.KeyResponse
	err := json.Unmarshal(res, &keyRes)

	if err != nil {
		return nil, JSONDecodeError
	}

	return &keyRes, nil
}
