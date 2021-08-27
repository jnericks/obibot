package telegram

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jnericks/obibot/internal/log"
)

type DeleteMessageParams struct {
	ChatID    int64 `json:"chat_id" validate:"required"`
	MessageID int64 `json:"message_id" validate:"required"`
}

func (c *client) DeleteMessage(ctx context.Context, params DeleteMessageParams) error {
	log.WithFields(ctx, log.Fields{
		"chatId":    params.ChatID,
		"messageId": params.MessageID,
	}).Info("deleting message")
	if err := c.validate.Struct(params); err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPost, c.url("deleteMessage"), params)
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
			log.WithError(ctx, err).WithField("status", status).Warn("deleteMessage failed")
		} else {
			log.WithFields(ctx, log.Fields{
				"status": status,
				"body":   string(b),
			}).Warn("deleteMessage failed")
		}
		return fmt.Errorf("deleteMessage failed: %s", status)
	}

	return nil
}
