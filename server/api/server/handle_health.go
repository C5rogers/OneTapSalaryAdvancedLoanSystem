package server

import (
	"fmt"
	"net/http"
)

func (s *Server) HandleServerHealtz(w http.ResponseWriter, r *http.Request) error {

	fmt.Fprintf(w, fmt.Sprintf("{\"status\": \"okay\"}"))

	return nil
}
