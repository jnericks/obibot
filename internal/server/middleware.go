package server

import (
	"net/http"

	"github.com/jnericks/obibot/internal/log"
)

func (s *server) middlewareInjectTraceID(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = log.WithTraceID(ctx)
		fn(w, r.WithContext(ctx))
	}
}
