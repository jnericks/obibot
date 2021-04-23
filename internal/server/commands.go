package server

import "github.com/jnericks/obibot/internal/commands"

func (s *server) commands() error {
	getStockPrice := commands.GetStockPrice(s.rapid)
	if err := s.manager.Register("/p", getStockPrice); err != nil {
		return err
	}
	if err := s.manager.Register("/price", getStockPrice); err != nil {
		return err
	}

	return nil
}
