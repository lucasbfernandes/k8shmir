package server

import (
	"io/ioutil"
	"k8s-smr/internal/models"
	"net/http"
)

func (s *Server) persistRequest(httpRequest *http.Request, requestId string) (*models.Request, error) {
	request, err := s.buildRequestObject(httpRequest, requestId)
	if err != nil {
		return nil, err
	}

	err = s.db.AppendRequest(request)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func (s *Server) buildRequestObject(nativeRequest *http.Request, requestId string) (*models.Request, error) {
	parsedBody, err := ioutil.ReadAll(nativeRequest.Body)
	if err != nil {
		return nil, err
	}

	return &models.Request{
		Id: requestId,
		RequestURI: nativeRequest.RequestURI,
		Host: nativeRequest.Host,
		Method: nativeRequest.Method,
		Url: nativeRequest.URL.String(),
		Headers: nativeRequest.Header,
		Body: parsedBody,
	}, nil
}
