package main

import (
	"req-smr/internal/api"
	"req-smr/internal/usecases"
)

func main() {
	usecases.WatchRequests()
	api.StartAPI()
}
