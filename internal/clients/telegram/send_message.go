package telegram

import (
	"context"
	"fmt"
	"io/ioutil"
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
		status := fmt.Sprintf("%d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
		if b, err := ioutil.ReadAll(resp.Body); err != nil {
			log.WithError(ctx, err).WithField("status", status).Warn("sendMessage failed")
		} else {
			log.WithFields(ctx, log.Fields{
				"status": status,
				"body":   string(b),
			}).Warn("sendMessage failed")
		}
		return fmt.Errorf("sendMessage failed: %s", status)
	}

	return nil
}
