package database

import (
	"context"
	"github.com/atomix/go-client/pkg/client"
	"k8s-smr/internal/config"
)

type RaftDatabase struct {
	client *client.Database

	primitiveName string
}

func New() (*RaftDatabase, error) {
	dbClient, err := getRaftDBClient()
	if err != nil {
		return nil, err
	}

	return &RaftDatabase{
		client: dbClient,
	}, nil
}

func getRaftDBClient() (*client.Database, error) {
	atomixControllerAddr := config.GetAtomixControllerAddress()
	atomixDBName := config.GetAtomixDBName()

	atomixClient, err := client.New(atomixControllerAddr)
	if err != nil {
		return nil, err
	}

	dbClient, err := atomixClient.GetDatabase(context.TODO(), atomixDBName)
	if err != nil {
		return nil, err
	}

	return dbClient, nil
}
