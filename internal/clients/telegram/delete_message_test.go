package telegram_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-playground/validator"
	"github.com/jnericks/obibot/internal/clients/telegram"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteMessageParamValidation(t *testing.T) {
	validParams := func() telegram.DeleteMessageParams {
		return telegram.DeleteMessageParams{
			ChatID:    123,
			MessageID: 456,
		}
	}

	tests := []struct {
		title    string
		field    string
		scenario func(params *telegram.DeleteMessageParams)
	}{
		{
			title: "ChatID is Zero",
			field: "ChatID",
			scenario: func(params *telegram.DeleteMessageParams) {
				params.ChatID = 0
			},
		},
		{
			title: "MessageID is Zero",
			field: "MessageID",
			scenario: func(params *telegram.DeleteMessageParams) {
				params.MessageID = 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			p := validParams()
			tt.scenario(&p)

			c, err := telegram.NewClient(http.DefaultClient, "fake-bot-token")
			require.NoError(t, err)
			err = c.DeleteMessage(context.Background(), p)

			verr, ok := err.(validator.ValidationErrors)
			require.True(t, ok)
			assert.Len(t, verr, 1)
			assert.Equal(t, tt.field, verr[0].Field())
		})
	}
}
