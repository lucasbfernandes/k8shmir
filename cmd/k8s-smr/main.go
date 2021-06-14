package main

import (
	"k8s-smr/internal/config"
	"k8s-smr/internal/server"
	"log"
	"os"
)

func main() {
	port := config.GetProxyServerPort()
	healthPort := config.GetHealthProxyServerPort()

	proxyServer, err := server.New(port, healthPort)
	if err != nil {
		log.Printf("failed to create server: %s\n", err.Error())
		os.Exit(1)
	}

	err = proxyServer.Start()
	if err != nil {
		log.Printf("failed to start server: %s\n", err.Error())
		os.Exit(1)
	}
}
