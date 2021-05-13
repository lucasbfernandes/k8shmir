package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"req-smr/internal/models"
	"req-smr/internal/services"

	"github.com/atomix/go-client/pkg/client/log"
)

func WatchRequests() error {
	fmt.Println("STEP:WATCH_REQUESTS")
	db, err := services.GetRaftDatabaseInstance()
	if err != nil {
		fmt.Printf("ERROR:GET_DATABASE %s\n", err)
		return err
	}

	fmt.Printf("STEP:GET_LOG\n")
	logPrimitive, err := db.GetLog(context.TODO(), "request-logs")
	if err != nil {
		fmt.Printf("ERROR:GET_LOG_REFERENCE %s\n", err)
		return err
	}

	fmt.Printf("STEP:WATCH_LOG_CHANNEL\n")
	channel := make(chan *log.Event)
	err = logPrimitive.Watch(context.Background(), channel)
	if err != nil {
		fmt.Printf("ERROR:WATCH_LOG %s\n", err)
		return err
	}

	go func() {
		for {
			fmt.Printf("STEP:WAITING_LOG_EVENT - channel: %s logPrimitive: %s\n", channel, logPrimitive)
			event := <-channel

			fmt.Printf("STEP:BYTE_ARRAY_TO_REQUEST %s\n", event)
			request, err := byteArrayToRequest(event.Entry.Value)
			if err != nil {
				fmt.Printf("ERROR:RECONSTRUCT_REQUEST %s\n", err)
				continue
			}

			fmt.Printf("STEP:CHECK_REQUEST_ORIGIN\n")
			if _, isPresent := services.RequestChanMap[request.Id]; isPresent {
				fmt.Printf("STEP:SAME_PROXY_REQUEST Id: %s\n", request.Id)
				services.RequestChanMap[request.Id] <- true
			} else {
				fmt.Printf("STEP:DIFFERENT_PROXY_REQUEST\n")
				_, err = services.ForwardRequest(request)
				if err != nil {
					fmt.Printf("ERROR:FORWARD_REQUEST %s\n", err)
					continue
				}
				fmt.Printf("STEP:FORWARD_REQUEST_SUCCEEDED %s\n", request)
			}
		}
	}()

	return nil
}

func byteArrayToRequest(serializedRequest []byte) (*models.Request, error) {
	var request models.Request
	err := json.Unmarshal(serializedRequest, &request)
	if err != nil {
		return nil, err
	}
	return &request, nil
}