package sentinel

import (
	"log"
	"net/http"
	"os"
	"strings"
)

type SentinelClient struct {
	user       string
	password   string
	httpClient *http.Client
	searchURL  string
	rows       int
}

func NewClient() *SentinelClient {
	credentials := strings.Split(os.Getenv("SENTINEL_CREDENTIALS"), ":")
	if len(credentials) < 2 {
		log.Fatalf("Please provide Sentinel credentials!")
	}

	return &SentinelClient{
		user:       credentials[0],
		password:   credentials[1],
		httpClient: &http.Client{},
		searchURL:  "https://scihub.copernicus.eu/dhus/search?q=",
		rows:       100,
	}
}
