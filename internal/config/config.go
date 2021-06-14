package config

import "os"

const (
	defaultProxyPort = "8080"

	defaultHealthProxyPort = "8081"

	defaultAppPort = "3000"

	defaultAtomixDBName = "raft-database"

	defaultAtomixControllerAddr = "atomix-controller.default.svc.cluster.local:5679"

	defaultAtomixLogPrimitiveName = "request-logs"
)

func GetProxyServerPort() string {
	port, envExists := os.LookupEnv("PROXY_PORT")
	if envExists && port != "" {
		return port
	}
	return defaultProxyPort
}

func GetHealthProxyServerPort() string {
	port, envExists := os.LookupEnv("PROXY_HEALTH_PORT")
	if envExists && port != "" {
		return port
	}
	return defaultHealthProxyPort
}

func GetApplicationPort() string {
	port, envExists := os.LookupEnv("APPLICATION_PORT")
	if envExists && port != "" {
		return port
	}
	return defaultAppPort
}

func GetAtomixDBName() string {
	dbName, envExists := os.LookupEnv("ATOMIX_DB_NAME")
	if envExists && dbName != "" {
		return dbName
	}
	return defaultAtomixDBName
}

func GetAtomixControllerAddress() string {
	addr, envExists := os.LookupEnv("ATOMIX_CONTROLLER_ADDRESS")
	if envExists && addr != "" {
		return addr
	}
	return defaultAtomixControllerAddr
}

func GetAtomixLogPrimitiveName() string {
	addr, envExists := os.LookupEnv("ATOMIX_LOG_PRIMITIVE_NAME")
	if envExists && addr != "" {
		return addr
	}
	return defaultAtomixLogPrimitiveName
}
