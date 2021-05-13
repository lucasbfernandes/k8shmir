package api

import (
	"fmt"
	"net/http"
	"os"
	"req-smr/internal/usecases"
)

var ProxyPort = os.Getenv("PROXY_PORT")

func StartAPI() {
	fmt.Printf("STEP:START_API_PORT: %s\n", ProxyPort)
	http.ListenAndServe(fmt.Sprintf(":%s", ProxyPort), &usecases.Proxy{})
}