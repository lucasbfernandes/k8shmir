package server

import (
	"bytes"
	"fmt"
	"k8s-smr/internal/config"
	"k8s-smr/internal/models"
	"net/http"
)

func (s *Server) forwardRequest(request *models.Request) (*http.Response, error) {
	proxyRequest, err := s.createHTTPRequestFromModel(request)
	if err != nil {
		return nil, err
	}

	httpClient := http.Client{}
	res, err := httpClient.Do(proxyRequest)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Server) createHTTPRequestFromModel(request *models.Request) (*http.Request, error) {
	applicationPort := config.GetApplicationPort()

	// We are always forwarding to 127.0.0.1 because both proxy and application reside on the same pod
	requestURL := fmt.Sprintf("http://127.0.0.1:%s%s", applicationPort, request.RequestURI)

	proxyRequest, err := http.NewRequest(request.Method, requestURL, bytes.NewBuffer(request.Body))
	if err != nil {
		return nil, err
	}
	proxyRequest.Header = request.Headers

	return proxyRequest, nil
}
