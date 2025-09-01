package server

import (
	"log/slog"
	"net/http"
)

func (s *Server) ApplyRoutes(mux *http.ServeMux) {

	mux.HandleFunc("GET /healthz", MakeAPI(s.HandleServerHealtz))

	mux.HandleFunc("POST /auth/login", MakeAPI(s.HandleLogin))

	mux.HandleFunc("POST /auth/register", MakeAPI(s.HandleRegister))

	mux.HandleFunc("POST /api/file_upload", MakeAPI(s.HandleFileUpload))

}

func MakeAPI(h func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := h(w, r); err != nil {
			slog.Error("internal server error", "err", err, "path", r.URL.Path)
		}
	}
}
