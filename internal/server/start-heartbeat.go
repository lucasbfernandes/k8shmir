package server

import (
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
	"k8s-smr/internal/config"
	"log"
	"net/http"
)

const (
	// TODO use environment variable
	getLastIndexPath = "/last-index"

	every5SecondsExpression = "@every 0h0m2s"
)

type LastIndexResponse struct {
	Index int `json:"index"`
}

func (s *Server) startHeartbeat() error {
	c := cron.New(cron.WithSeconds())

	_, err := c.AddFunc(every5SecondsExpression, s.doHeartbeat)
	if err != nil {
		return err
	}

	c.Start()
	return nil
}

func (s *Server) doHeartbeat() {
	lastIndex, err := s.getAppCurrentAppliedLogIndex()
	if err != nil {
		log.Printf("failed to get last applied index: %v\n", err)
		s.isSynced = false
	}

	log.Printf("last index response: %v\n", lastIndex)
}

func (s *Server) getAppCurrentAppliedLogIndex() (*LastIndexResponse, error) {
	var lastIndexResponse LastIndexResponse

	applicationPort := config.GetApplicationPort()

	// We are always requesting to 127.0.0.1 because both proxy and application reside on the same pod
	requestURL := fmt.Sprintf("http://127.0.0.1:%s%s", applicationPort, getLastIndexPath)

	proxyRequest, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	httpClient := http.Client{}
	res, err := httpClient.Do(proxyRequest)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(res.Body).Decode(&lastIndexResponse)
	if err != nil {
		return nil, err
	}

	return &lastIndexResponse, nil
}
