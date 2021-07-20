package server

import (
	"io/ioutil"
	"net/http"

	"github.com/jnericks/obibot/internal/log"
)

func (s *server) injectTraceID(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = log.WithTraceID(ctx)
		fn(w, r.WithContext(ctx))
	}
}

func (s *server) httpMethod(fn http.HandlerFunc, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rm := r.Method
		if rm == "" {
			rm = http.MethodGet
		}

		if rm != method {
			b, _ := ioutil.ReadAll(r.Body)
			log.WithFields(r.Context(), log.Fields{
				"method": r.Method,
				"url":    r.URL.String(),
				"body":   string(b),
			}).Warn("invalid request method")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		fn(w, r)
	}
}
