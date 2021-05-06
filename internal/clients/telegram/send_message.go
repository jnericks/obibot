package telegram

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jnericks/obibot/internal/log"
)

const (
	ParseModeMarkdown = "MarkdownV2"
)

type SendMessageParams struct {
	ChatID    int64  `json:"chat_id" validate:"gt=1"`
	Text      string `json:"text" validate:"required"`
	ParseMode string `json:"parse_mode,omitempty" validate:"omitempty,oneof=MarkdownV2"`
}

func (c *client) SendMessage(ctx context.Context, params SendMessageParams) error {
	log.WithFields(ctx, log.Fields{
		"chatId":    params.ChatID,
		"text":      params.Text,
		"parseMode": params.ParseMode,
	}).Info("sending message")
	if err := c.validate.Struct(params); err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPost, c.url("sendMessage"), params)
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
