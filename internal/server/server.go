package server

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"k8s-smr/internal/database"
	"log"
	"net/http"
)

type Server struct {
	port string

	healthPort string

	db *database.RaftDatabase

	incomingRequestsMap map[string]bool

	isSynced bool
}

// TODO use context instead of a request map?
func New(port string, healthPort string) (*Server, error) {
	raftDatabase, err := database.New()
	if err != nil {
		return nil, err
	}

	return &Server{
		port: port,
		healthPort: healthPort,
		db: raftDatabase,
		incomingRequestsMap: make(map[string]bool),
	}, nil
}

// TODO add support for TLS?
func (s *Server) Start() error {
	err := s.watchRequests()
	if err != nil {
		return err
	}

	s.isSynced = false

	err = s.startHeartbeat()
	if err != nil {
		return err
	}

	go s.startHealthCheckServer()

	err = http.ListenAndServe(fmt.Sprintf(":%s", s.port), s)
	if err != nil {
		return err
	}

	return nil
}

// TODO create dto and return it instead of model?
func (s *Server) ServeHTTP(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	requestId := uuid.New().String()
	s.incomingRequestsMap[requestId] = true

	request, err := s.buildRequestObject(httpRequest, requestId)
	if err != nil {
		log.Printf("failed to create request object: %s\n", err)
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	logEntry, err := s.persistRequest(request)
	if err != nil {
		log.Printf("failed to persist request: %s\n", err)
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO improve error handling - might add inconsistency
	res, err := s.forwardRequest(request, logEntry)
	if err != nil {
		log.Printf("failed to forward request: %s\n", err)
		http.Error(responseWriter, err.Error(), http.StatusBadGateway)
		return
	}

	log.Printf("forwarded request with application response: %+v\n", res)

	err = s.writeResponse(responseWriter, res)
	if err != nil {
		log.Printf("failed to write response back to client: %s\n", err)
		http.Error(responseWriter, err.Error(), http.StatusBadGateway)
		return
	}

	err = res.Body.Close()
	if err != nil {
		log.Printf("failed to close response body: %s\n", err)
	}
}

func (s *Server) writeResponse(responseWriter http.ResponseWriter, res *http.Response) error {
	for name, values := range res.Header {
		responseWriter.Header()[name] = values
	}

	responseWriter.WriteHeader(res.StatusCode)

	_, err := io.Copy(responseWriter, res.Body)
	if err != nil {
		return err
	}

	return nil
}
