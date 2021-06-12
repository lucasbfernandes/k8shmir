package database

import (
	"context"
	"encoding/json"
	"github.com/atomix/go-client/pkg/client/log"
	"k8s-smr/internal/config"
	"k8s-smr/internal/models"
)

func (db *RaftDatabase) AppendRequest(request *models.Request) (*log.Entry, error) {
	logPrimitiveName := config.GetAtomixLogPrimitiveName()

	logPrimitive, err := db.client.GetLog(context.TODO(), logPrimitiveName)
	if err != nil {
		return nil, err
	}

	serializedRequest, err := db.requestToByteArray(request)
	if err != nil {
		return nil, err
	}

	entry, err := logPrimitive.Append(context.TODO(), serializedRequest)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

func (db *RaftDatabase) GetRequestsWatchChannel() (chan *log.Event, error) {
	logPrimitiveName := config.GetAtomixLogPrimitiveName()

	logPrimitive, err := db.client.GetLog(context.TODO(), logPrimitiveName)
	if err != nil {
		return nil, err
	}

	requestsChan := make(chan *log.Event)
	err = logPrimitive.Watch(context.TODO(), requestsChan)
	if err != nil {
		return nil, err
	}

	return requestsChan, nil
}

func (db *RaftDatabase) requestToByteArray(request *models.Request) ([]byte, error) {
	serializedRequest, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	return serializedRequest, nil
}
