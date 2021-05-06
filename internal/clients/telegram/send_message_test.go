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

func TestSendMessageParamValidation(t *testing.T) {
	validParams := func() telegram.SendMessageParams {
		return telegram.SendMessageParams{
			ChatID:    123,
			Text:      "something",
			ParseMode: "",
		}
	}

	tests := []struct {
		title    string
		field    string
		scenario func(params *telegram.SendMessageParams)
	}{
		{
			title: "ChatID is Zero",
			field: "ChatID",
			scenario: func(params *telegram.SendMessageParams) {
				params.ChatID = 0
			},
		},
		{
			title: "ChatID is Negative",
			field: "ChatID",
			scenario: func(params *telegram.SendMessageParams) {
				params.ChatID = -1
			},
		},
		{
			title: "Text is empty",
			field: "Text",
			scenario: func(params *telegram.SendMessageParams) {
				params.Text = ""
			},
		},
		{
			title: "ParseMode is not valid",
			field: "ParseMode",
			scenario: func(params *telegram.SendMessageParams) {
				params.ParseMode = "INVALID"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			p := validParams()
			tt.scenario(&p)

			c, err := telegram.NewClient(http.DefaultClient, "fake-bot-token")
			require.NoError(t, err)
			err = c.SendMessage(context.Background(), p)

			verr, ok := err.(validator.ValidationErrors)
			require.True(t, ok)
			assert.Len(t, verr, 1)
			assert.Equal(t, tt.field, verr[0].Field())
		})
	}
}
