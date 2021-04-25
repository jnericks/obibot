package server

import "github.com/jnericks/obibot/internal/commands"

func (s *server) commands() error {
	if err := s.manager.Register("!p", commands.GetSymbolQuote(s.iex)); err != nil {
		return err
	}

	return nil
}
