package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"net/http"
	"req-smr/internal/models"
	"req-smr/internal/services"
)

type Proxy struct{}

func (proxy *Proxy) ServeHTTP(responseWriter http.ResponseWriter, httpRequest *http.Request) {

	requestId := uuid.New().String()
	services.RequestChanMap[requestId] = make(chan bool)

	fmt.Printf("STEP:INCOMING_REQUEST %s\n", httpRequest)
	fmt.Println("STEP:BUILD_REQUEST_OBJECT")
	request, err := buildRequestObject(httpRequest, requestId)
	if err != nil {
		fmt.Printf("ERROR:BUILD_REQUEST_OBJECT %s\n", err)
		http.Error(responseWriter, err.Error(), http.StatusBadGateway)
		return
	}

	fmt.Println("STEP:PERSIST_LOG")
	err = persistRequest(request)
	if err != nil {
		fmt.Printf("ERROR:PERSIST_LOG %s\n", err)
		// Return error message for client
		http.Error(responseWriter, err.Error(), http.StatusBadGateway)
		return
	}

	fmt.Printf("STEP:WAITING_REQUEST_CHANNEL\n")
	<-services.RequestChanMap[requestId]

	fmt.Println("STEP:PROXY_HTTP_FORWARD_REQUEST")
	res, err := services.ForwardRequest(request)
	if err != nil {
		fmt.Printf("ERROR:FORWARD_REQUEST %s\n", err)
		http.Error(responseWriter, err.Error(), http.StatusBadGateway)
		return
	}

	fmt.Println("STEP:PROXY_HTTP_WRITE_RESPONSE")
	writeResponse(responseWriter, res)
}

func persistRequest(request *models.Request) error {
	fmt.Println("START:SET_REQUEST")
	db, err := services.GetRaftDatabaseInstance()
	if err != nil {
		fmt.Printf("ERROR:GET_DATABASE %s\n", err)
		return err
	}

	log, err := db.GetLog(context.TODO(), "request-logs")
	if err != nil {
		fmt.Printf("ERROR:GET_LOG_REFERENCE %s\n", err)
		return err
	}

	serializedRequest, err := requestToByteArray(request)
	if err != nil {
		fmt.Printf("ERROR:SERIALIZE_REQUEST %s\n", err)
		return err
	}
	fmt.Printf("GET:SERIALIZED_REQUEST %s\n", serializedRequest)

	_, err = log.Append(context.TODO(), serializedRequest)
	if err != nil {
		fmt.Printf("ERROR:APPEND_LOG %s\n", err)
		return err
	}

	fmt.Println("FINISH:SET_APPEND")
	return nil
}

func requestToByteArray(request *models.Request) ([]byte, error) {
	serializedRequest, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	return serializedRequest, nil
}

func buildRequestObject(nativeRequest *http.Request, requestId string) (*models.Request, error) {
	parsedBody, err := parseRequestBody(nativeRequest.Body)
	if err != nil {
		return nil, err
	}

	request := &models.Request{
		Id: requestId,
		RequestURI: nativeRequest.RequestURI,
		Host: nativeRequest.Host,
		Method: nativeRequest.Method,
		Url: nativeRequest.URL.String(),
		Headers: nativeRequest.Header,
		Body: parsedBody,
	}
	return request, nil
}

func parseRequestBody(requestBody io.ReadCloser) ([]byte, error) {
	body, err := ioutil.ReadAll(requestBody)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func writeResponse(responseWriter http.ResponseWriter, res *http.Response) {
	fmt.Println("STEP:WRITE_RESPONSE")
	for name, values := range res.Header {
		responseWriter.Header()[name] = values
	}
	responseWriter.Header().Set("Server", "req-smr")
	responseWriter.WriteHeader(res.StatusCode)
	io.Copy(responseWriter, res.Body)
	res.Body.Close()
	fmt.Println("STEP:WRITE_RESPONSE_SUCCEEDED")
}