package sentinel

import (
	"log"
	"os"
	"strings"
)

type SentinelClient struct {
	user     string
	password string
}

func NewClient() *SentinelClient {
	credentials := strings.Split(os.Getenv("SENTINEL_CREDENTIALS"), ":")
	if len(credentials) < 2 {
		log.Fatalf("Please provide Sentinel credentials!")
	}

	return &SentinelClient{
		user:     credentials[0],
		password: credentials[1],
	}
}
