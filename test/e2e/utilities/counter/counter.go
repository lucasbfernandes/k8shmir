package counter

import (
	"github.com/go-resty/resty/v2"
	"os"
	"os/exec"

	"errors"
	"fmt"
)

const (
	URL1 = "http://localhost/counter1"

	URL2 = "http://localhost/counter2"

	IncOP = "INC"

	DecOP = "DEC"
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

func DoResetCounter() error {
	client := resty.New()
	response, err := client.R().
		Post(fmt.Sprintf("%s/integer/reset", URL1))

	if err != nil {
		return err
	}

	if response.IsError() {
		return errors.New(fmt.Sprintf("failed with status code: %d", response.StatusCode()))
	}

	return nil
}

func DoAlternateRequest(index int, url string, incVal int, decVal int) error {
	if index % 2 == 0 {
		err := DoPostCounterRequest(url, IncOP, incVal)
		return err
	} else {
		err := DoPostCounterRequest(url, DecOP, decVal)
		return err
	}
}

func ExecuteAndWaitScriptFile(scriptPath string) error {
	execCmd := exec.Command("/bin/sh", scriptPath)

	execCmd.Stderr = os.Stderr
	execCmd.Stdout = os.Stdout

	err := execCmd.Start()
	if err != nil {
		return err
	}

	err = execCmd.Wait()
	if err != nil {
		return err
	}

	return nil
}
