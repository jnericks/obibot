package server

func (s *server) routes() {
	s.router.HandleFunc("/api", s.middlewareInjectTraceID(s.handleAPI()))
	s.router.HandleFunc("/telegram", s.middlewareInjectTraceID(s.handleTelegram()))
}
