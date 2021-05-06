package server

import (
	"net/http"
	"strings"

	"github.com/jnericks/obibot/internal/clients/telegram"
	"github.com/jnericks/obibot/internal/commands"
	"github.com/jnericks/obibot/internal/log"
)

func (s *server) handleAPI() http.HandlerFunc {
	type Request struct {
		Command string `json:"command"`
		Args    string `json:"args"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var data Request
		if err := s.decode(w, r, &data); err != nil {
			log.WithError(ctx, err).Warn("decoding api request body")
			s.respond(w, r, nil, http.StatusBadRequest)
			return
		}

		out, err := s.manager.Exec(r.Context(), data.Command, commands.Input{
			Args: strings.Split(data.Args, " "),
		})
		if err != nil || out == nil {
			log.WithError(ctx, err).Warn()
			s.respond(w, r, nil, http.StatusBadRequest)
		}

		params := telegram.SendMessageParams{
			ChatID: 1,
			Text:   out.Response,
		}
		if out.Markdown {
			params.ParseMode = "MarkdownV2"
		}
		s.respond(w, r, params, http.StatusOK)
	}
}
