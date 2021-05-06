package server

import (
	"net/http"
	"strings"

	"github.com/jnericks/obibot/internal/clients/telegram"
	"github.com/jnericks/obibot/internal/commands"
	"github.com/jnericks/obibot/internal/log"
)

func (s *server) handleTelegram() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var data telegram.Request
		if err := s.decode(w, r, &data); err != nil {
			log.WithError(ctx, err).Warn("decoding telegram request body")
			s.respond(w, r, nil, http.StatusBadRequest)
			return
		}

		in := strings.Split(data.Message.Text, " ")
		if len(in) < 2 {
			return
		}

		cmd := in[0]
		cmdInput := commands.Input{Args: in[1:]}
		if !s.manager.CanExec(cmd) {
			// do nothing since chat message is just regular text, not a command
			return
		}

		log.WithFields(ctx, log.Fields{
			"data":    data,
			"command": cmd,
			"input":   cmdInput,
		}).Info("handling telegram request")

		out, err := s.manager.Exec(ctx, cmd, cmdInput)
		if err != nil {
			log.WithError(ctx, err).WithFields(log.Fields{
				"command": cmd,
				"input":   cmdInput,
			}).Warn()
			s.respond(w, r, nil, http.StatusBadRequest)
			return
		}

		params := telegram.SendMessageParams{
			ChatID: data.Message.Chat.ID,
			Text:   out.Response,
		}
		if out.Markdown {
			params.ParseMode = telegram.ParseModeMarkdown
		}
		if err := s.telegram.SendMessage(ctx, params); err != nil {
			log.WithError(ctx, err).Error("sending telegram message")
			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
