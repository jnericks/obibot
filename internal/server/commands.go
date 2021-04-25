package server

import "github.com/jnericks/obibot/internal/commands"

func (s *server) commands() error {
	for cmd, fn := range map[string]commands.Func{
		"/s": commands.GetStock(s.iex),
		"/c": commands.GetCrypto(s.iex),
	} {
		if err := s.manager.Register(cmd, fn); err != nil {
			return err
		}
	}
	return nil
}
