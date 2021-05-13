package services

import (
	"context"
	"fmt"
	"sync"

	"github.com/atomix/go-client/pkg/client"
)

var dbInstance *client.Database
var syncError error
var once sync.Once

func GetRaftDatabaseInstance() (*client.Database, error) {
	fmt.Println("STEP:GET_ATOMIX_DB")
	once.Do(func() {
		fmt.Println("STEP:CREATE_ATOMIX_CONNECTION")
		atomix, err := client.New("atomix-controller.default.svc.cluster.local:5679")
		if err != nil {
			fmt.Printf("ERROR:GET_ATOMIX_CONNECTION %s\n", err)
			syncError = err
		}

		fmt.Println("STEP:GET_ATOMIX_DB_INSTANCE")
		dbInstance, err = atomix.GetDatabase(context.TODO(), "raft-database")
		if err != nil {
			fmt.Printf("ERROR:GET_DATABASE_CONNECTION %s\n", err)
			syncError = err
		}
	})

	fmt.Printf("STEP:GET_DB_SUCCESS %+v\n", dbInstance)
	return dbInstance, syncError
}
