package server

import (
	"errors"
	"net/http"
	"strings"

	"github.com/jnericks/obibot/internal/clients/telegram"
	"github.com/jnericks/obibot/internal/commands"
	log "github.com/sirupsen/logrus"
)

func (s *server) routes() {
	s.router.HandleFunc("/api", s.handleAPI())
	s.router.HandleFunc("/", s.middlewareExample(s.handleTelegram()))
}

func (s *server) handleAPI() http.HandlerFunc {
	type Request struct {
		Command string `json:"command"`
		Args    string `json:"args"`
	}
	type Response struct {
		Message string `json:"message"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var data Request
		if err := s.decode(w, r, &data); err != nil {
			log.WithError(err).Warn("decoding api request body")
			s.respond(w, r, nil, http.StatusBadRequest)
			return
		}

		out, err := s.manager.Exec(r.Context(), data.Command, commands.Input{
			Args: strings.Split(data.Args, " "),
		})
		if err != nil {
			log.WithError(err).Warn()
			s.respond(w, r, nil, http.StatusBadRequest)
		}

		s.respond(w, r, Response{Message: out.Response}, http.StatusOK)
	}
}

func (s *server) handleTelegram() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data telegram.Request
		if err := s.decode(w, r, &data); err != nil {
			log.WithError(err).Warn("decoding telegram request body")
			s.respond(w, r, nil, http.StatusBadRequest)
			return
		}

		in := strings.Split(data.Message.Text, " ")
		if len(in) < 2 {
			return
		}

		out, err := s.manager.Exec(r.Context(), in[0], commands.Input{Args: in[1:]})
		if err != nil {
			if errors.Is(err, commands.ErrNotSupported{}) {
				// do nothing since chat message is just regular text, not a command
				return
			}

			log.WithError(err).WithField("input", in).Warn()
			s.respond(w, r, nil, http.StatusBadRequest)
			return
		}

		if err := s.telegram.SendMessage(r.Context(), telegram.SendMessageParams{
			ChatID: data.Message.Chat.ID,
			Text:   out.Response,
		}); err != nil {
			log.WithError(err).Error("sending telegram message")
			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
