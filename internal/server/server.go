package server

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/jnericks/obibot/internal/clients/rapid"
	"github.com/jnericks/obibot/internal/clients/telegram"
	"github.com/jnericks/obibot/internal/commands"
)

type Config struct{}

type Dependencies struct {
	Telegram telegram.Client `validate:"required"`
	Rapid    rapid.Client    `validate:"required"`
}

type server struct {
	router  http.ServeMux
	manager commands.Manager

	telegram telegram.Client
	rapid    rapid.Client
}

func NewServer(config Config, deps Dependencies) (http.Handler, error) {
	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return nil, err
	}
	if err := validate.Struct(deps); err != nil {
		return nil, err
	}

	s := &server{
		router:  http.ServeMux{},
		manager: commands.NewManager(),

		rapid:    deps.Rapid,
		telegram: deps.Telegram,
	}
	s.routes()
	if err := s.commands(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
