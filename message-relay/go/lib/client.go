package lib

import (
	"errors"

	"github.com/Scotiacon-Tech/libs/message-relay/go/requests"
)

type Client struct {
	Key string
}

func NewClient() *Client {
	return &Client{}
}

func (client Client) NewSendRequest() *requests.SendRequest {
	return &requests.SendRequest{
		From:    "",
		To:      "",
		Subject: "",
		Body:    "",
	}
}

func (client Client) SendMessage(storedKey string, service string, req *requests.SendRequest) (*SendResult, error) {
	if service == "" {
		return nil, errors.New("Service name required")
	}

	key := storedKey

	sendRes, sendErr := client.RequestSend(key, service, req)

	if sendErr == KeyInvalidError {
		jwtRes, jwtErr := client.RequestJWT()

		if jwtErr != nil {
			return nil, errors.New("Failed to fetch JWT")
		}

		keyRes, keyErr := client.RequestKey(jwtRes.AccessToken)

		if keyErr != nil {
			return nil, errors.New("Failed to fetch key")
		}

		key = keyRes.Key

		sendRes, sendErr = client.RequestSend(key, service, req)
	}

	if sendErr != nil {
		return nil, errors.New("Error sending message")
	}

	res := SendResult{
		NewKey:    key,
		MessageID: sendRes.MessageID,
	}

	return &res, nil
}
