package services

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"req-smr/internal/models"
)

var RequestChanMap = make(map[string]chan bool)

var ApplicationPort = os.Getenv("SERVICE_PORT")

func ForwardRequest(request *models.Request) (*http.Response, error) {
	proxyUrl := fmt.Sprintf("http://127.0.0.1:%s%s", ApplicationPort, request.RequestURI)
	httpClient := http.Client{}
	proxyReq, err := http.NewRequest(request.Method, proxyUrl, bytes.NewBuffer(request.Body))
	if err != nil {
		fmt.Printf("ERROR:CREATE_NEW_REQUEST_OBJECT %s\n", err)
		return nil, err
	}

	proxyReq.Header = request.Headers

	fmt.Printf("STEP:DO_REQUEST %s\n", proxyReq)
	res, err := httpClient.Do(proxyReq)
	if err != nil {
		fmt.Printf("ERROR:HTTP_CLIENT_DO %s\n", err)
		return nil, err
	}
	return res, nil
}