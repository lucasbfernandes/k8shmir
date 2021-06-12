package server

import (
	"github.com/atomix/go-client/pkg/client/log"
	"io/ioutil"
	"k8s-smr/internal/models"
	"net/http"
)

func (s *Server) persistRequest(request *models.Request) (*log.Entry, error) {
	entry, err := s.db.AppendRequest(request)
	if err != nil {
		return nil, err
	}

	return entry, nil
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
