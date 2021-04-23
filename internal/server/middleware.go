package server

import "net/http"

func (s *server) middlewareExample(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// do pre handler things
		fn(w, r)
		// do post handler things
	}
}
