package telegram

import (
	"context"
	"fmt"
	"net/http"
)

type SendMessageParams struct {
	ChatID    int64  `json:"chat_id" validate:"required"`
	Text      string `json:"text" validate:"required"`
	ParseMode string `json:"parse_mode,omitempty" validate:"oneof=MarkdownV2"`
}

func (c *client) SendMessage(ctx context.Context, params SendMessageParams) error {
	if err := c.validate.Struct(params); err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPost, c.url()+endpointSendMessage, params)
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("sendMessage failed: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	return nil
}
