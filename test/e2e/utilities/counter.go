package e2e_test_utilities

import (
	"github.com/go-resty/resty/v2"

	"errors"
	"fmt"
)

const (
	Counter1URL = "http://localhost/counter1"

	Counter2URL = "http://localhost/counter2"
)

type CounterRequest struct {
	Op string	`json:"op"`
	Value int	`json:"value"`
}

type CounterResponse struct {
	Value int	`json:"value"`
}

func DoPostCounterRequest(url string, op string, value int) error {
	requestObj := CounterRequest{
		Op: op,
		Value: value,
	}

	client := resty.New()
	response, err := client.R().
		SetBody(requestObj).
		Post(fmt.Sprintf("%s/integer", url))

	if err != nil {
		return err
	}

	if response.IsError() {
		return errors.New(fmt.Sprintf("failed with status code: %d", response.StatusCode()))
	}

	return nil
}

func DoGetCounterRequest(url string) (*CounterResponse, error) {
	var counterResponse CounterResponse

	client := resty.New()
	response, err := client.R().
		SetResult(&counterResponse).
		Get(fmt.Sprintf("%s/integer", url))

	if err != nil {
		return nil, err
	}

	if response.IsError() {
		return nil, errors.New(fmt.Sprintf("failed with status code: %d", response.StatusCode()))
	}

	return &counterResponse, nil
}
