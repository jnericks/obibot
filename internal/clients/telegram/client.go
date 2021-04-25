package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator"
)

const (
	endpointSendMessage = "/sendMessage"
)

type Client interface {
	SendMessage(context.Context, SendMessageParams) error
}

func NewClient(httpClient *http.Client, botToken string) (Client, error) {
	validate := validator.New()
	if err := validate.Struct(struct {
		BotToken string `validate:"required"`
	}{
		BotToken: botToken,
	}); err != nil {
		return nil, err
	}

	return &client{
		botToken: botToken,
		validate: validate,
		http:     httpClient,
	}, nil
}

type client struct {
	botToken string
	validate *validator.Validate
	http     *http.Client
}

func (c *client) url() string {
	return fmt.Sprintf("https://api.telegram.org/bot%s", c.botToken)
}

func (c *client) newRequest(ctx context.Context, method, url string, data interface{}) (*http.Request, error) {
	var body io.Reader
	if data != nil {
		b, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req.WithContext(ctx), nil
}