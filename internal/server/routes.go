package server

import (
	"io/ioutil"
	"net/http"

	"github.com/jnericks/obibot/internal/log"
)

func (s *server) routes() {
	s.route(http.MethodPost, "/api", s.handleAPI())
	s.route(http.MethodPost, "/telegram", s.handleTelegram())
	s.route(http.MethodGet, "/ping", s.healthCheck())
	s.routeCatchAll()
}

func (s *server) route(method, pattern string, fn http.HandlerFunc) {
	s.router.HandleFunc(pattern, s.injectTraceID(s.httpMethod(fn, method)))
}

func (s *server) routeCatchAll() {
	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b, _ := ioutil.ReadAll(r.Body)
		log.WithFields(log.ProcessContext, log.Fields{
			"method": r.Method,
			"url":    r.URL.String(),
			"body":   string(b),
		}).Warn("invalid request")
		w.WriteHeader(http.StatusNotFound)
	})
}

func (s *server) healthCheck() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("PONG"))
		if err != nil {
			log.WithError(r.Context(), err).Error()
		}
	}
}
