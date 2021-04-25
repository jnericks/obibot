package server

import "github.com/jnericks/obibot/internal/commands"

func (s *server) commands() error {
	if err := s.manager.Register("!s", commands.GetStockQuote(s.iex)); err != nil {
		return err
	}
	if err := s.manager.Register("!c", commands.GetCryptoPrice(s.iex)); err != nil {
		return err
	}

	return nil
}
