package server

import (
	"encoding/json"
	"fmt"
	atomixLog "github.com/atomix/go-client/pkg/client/log"
	"github.com/robfig/cron/v3"
	"k8s-smr/internal/config"
	"log"
	"net/http"
)

const (
	// TODO use environment variable
	getLastIndexPath = "/last-index"

	every5SecondsExpression = "@every 0h0m2s"
)

type LastIndexResponse struct {
	Index int `json:"index"`
}

func (s *Server) startHeartbeat() error {
	c := cron.New(
		cron.WithSeconds(),
		cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)),
	)

	_, err := c.AddFunc(every5SecondsExpression, s.doHeartbeat)
	if err != nil {
		return err
	}

	c.Start()
	return nil
}

func (s *Server) doHeartbeat() {
	lastIndex, err := s.getAppCurrentAppliedLogIndex()
	if err != nil {
		log.Printf("failed to get last applied index: %v\n", err)
		s.isSynced = false
		return
	}

	err = s.syncData(lastIndex)
	if err != nil {
		log.Printf("failed to sync data: %v\n", err)
		s.isSynced = false
		return
	}

	s.isSynced = true
	log.Printf("data synced successfuly\n")
}

func (s *Server) getAppCurrentAppliedLogIndex() (*LastIndexResponse, error) {
	var lastIndexResponse LastIndexResponse

	applicationPort := config.GetApplicationPort()

	// We are always requesting to 127.0.0.1 because both proxy and application reside on the same pod
	requestURL := fmt.Sprintf("http://127.0.0.1:%s%s", applicationPort, getLastIndexPath)

	proxyRequest, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	httpClient := http.Client{}
	res, err := httpClient.Do(proxyRequest)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(res.Body).Decode(&lastIndexResponse)
	if err != nil {
		return nil, err
	}

	return &lastIndexResponse, nil
}

func (s *Server) syncData(response *LastIndexResponse) error {
	lastAppliedIndex := uint64(response.Index)
	lastAtomixIndex, err := s.db.GetLastIndex()
	if err != nil {
		return err
	}

	err = s.applyRequestsDiff(lastAppliedIndex, uint64(*lastAtomixIndex))
	if err != nil {
		return err
	}

	err = s.clearWatchQueue(uint64(*lastAtomixIndex))
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) applyRequestsDiff(lastAppliedIndex uint64, lastAtomixIndex uint64) error {
	for i := lastAppliedIndex + 1; i <= lastAtomixIndex; i++ {
		entry, err := s.db.GetByIndex(atomixLog.Index(i))
		if err != nil {
			return err
		}

		request, err := s.byteArrayToRequest(entry.Value)
		if err != nil {
			return err
		}

		_, err = s.forwardRequest(request, entry)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) clearWatchQueue(lastAtomixIndex uint64) error {
	for len(s.watchQueue) > 0 {
		request := s.watchQueue[0].request
		entry := s.watchQueue[0].logEntry
		queueIndex := uint64(entry.Index)

		if queueIndex >= lastAtomixIndex {
			_, err := s.forwardRequest(request, entry)
			if err != nil {
				return err
			}
		}

		s.watchQueue = s.watchQueue[1:]
	}
	return nil
}
