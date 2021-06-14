package server

import (
	"fmt"
	"log"
	"net/http"
)

const (
	healthcheckPath = "/healthcheck"
)

func (s *Server) startHealthCheckServer() {
	go func() {
		http.HandleFunc(healthcheckPath, func(responseWriter http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				s.writeHealthcheckResponse(responseWriter)
			} else {
				http.Error(responseWriter, "invalid request method", http.StatusMethodNotAllowed)
			}
		})

		err := http.ListenAndServe(fmt.Sprintf(":%s", s.healthPort), nil)
		if err != nil {
			log.Printf("failed to start healthcheck server: %v\n", err)
		}
	}()
}

func (s *Server) writeHealthcheckResponse(responseWriter http.ResponseWriter) {
	if s.isSynced {
		responseWriter.WriteHeader(http.StatusOK)
		_, err := responseWriter.Write(nil)
		if err != nil {
			http.Error(responseWriter, "state not synced", http.StatusInternalServerError)
		}
	} else {
		// TODO review this status
		http.Error(responseWriter, "state not synced", http.StatusPreconditionFailed)
	}
}
