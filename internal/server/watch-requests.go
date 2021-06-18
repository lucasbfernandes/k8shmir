package server

import (
	"encoding/json"
	atomixLog "github.com/atomix/go-client/pkg/client/log"
	"k8s-smr/internal/models"
	"log"
)

func (s *Server) watchRequests() error {
	watchChan, err := s.db.GetRequestsWatchChannel()
	if err != nil {
		return err
	}

	go s.processObservedRequests(watchChan)

	return nil
}

func (s *Server) processObservedRequests(watchChan chan *atomixLog.Event) {
	for {
		event := <-watchChan

		// TODO improve error handling - might add inconsistency
		request, err := s.byteArrayToRequest(event.Entry.Value)
		if err != nil {
			log.Printf("failed to convert byte array to request: %s\n", err)
			continue
		}

		// TODO improve error handling - might add inconsistency
		if _, requestExists := s.incomingRequestsMap[request.Id]; !requestExists {
			if !s.isSynced {
				s.watchQueue = append(s.watchQueue, WatchQueueEntry{request, event.Entry})
				continue
			}

			_, err = s.forwardRequest(request, event.Entry)
			if err != nil {
				log.Printf("watch - failed to forward request: %s\n", err)
				s.isSynced = false
				continue
			}
		}
	}
}

func (s *Server) byteArrayToRequest(serializedRequest []byte) (*models.Request, error) {
	var request models.Request
	err := json.Unmarshal(serializedRequest, &request)
	if err != nil {
		return nil, err
	}
	return &request, nil
}
