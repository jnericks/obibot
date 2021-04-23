package server

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (s *server) decode(_ http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func (s *server) respond(w http.ResponseWriter, _ *http.Request, data interface{}, status int) {
	if data != nil {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.WithError(err).Error("encoding data to response")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	w.WriteHeader(status)
}
