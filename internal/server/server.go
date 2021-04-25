package server

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/jnericks/obibot/internal/clients/iex"
	"github.com/jnericks/obibot/internal/clients/telegram"
	"github.com/jnericks/obibot/internal/commands"
)

type Config struct{}

type Dependencies struct {
	Telegram telegram.Client `validate:"required"`
	IEX      iex.Client      `validate:"required"`
}

type server struct {
	router  http.ServeMux
	manager commands.Manager

	telegram telegram.Client
	iex      iex.Client
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

		telegram: deps.Telegram,
		iex:      deps.IEX,
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
