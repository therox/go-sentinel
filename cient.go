package sentinel

import (
	"log"
	"os"
	"strings"
)

type SentinelClient struct {
	user     string
	password string
	oData    struct {
		productsURL string
	}
	openAPI struct {
		searchURL string
		rows      int
	}
}

func NewClient() *SentinelClient {
	credentials := strings.Split(os.Getenv("SENTINEL_CREDENTIALS"), ":")
	if len(credentials) < 2 {
		log.Fatalf("Please provide Sentinel credentials!")
	}

	return &SentinelClient{
		user:     credentials[0],
		password: credentials[1],
		oData: struct {
			productsURL string
		}{
			productsURL: "https://scihub.copernicus.eu/dhus/odata/v1/Products",
		},
		openAPI: struct {
			searchURL string
			rows      int
		}{
			searchURL: "https://scihub.copernicus.eu/dhus/search",
			rows:      100,
		},
	}
}
