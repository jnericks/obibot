package server

import "github.com/jnericks/obibot/internal/commands"

func (s *server) commands() error {
	getStockQuote := commands.GetStockQuote(s.iex)
	if err := s.manager.Register("/s", getStockQuote); err != nil {
		return err
	}
	if err := s.manager.Register("!s", getStockQuote); err != nil {
		return err
	}

	getCryptoPrice := commands.GetCryptoPrice(s.iex)
	if err := s.manager.Register("/c", getCryptoPrice); err != nil {
		return err
	}
	if err := s.manager.Register("!c", getCryptoPrice); err != nil {
		return err
	}

	return nil
}
