package sentinel

import (
	"fmt"
	"net/http"
	"time"

	sentinel_engine "github.com/therox/go-sentinel/backend/sentinel"
)

type SentinelClient struct {
	user       string
	password   string
	httpClient *http.Client
	searchURL  string
	rows       int
	dlEngine   dlEngine
}

func NewClient(user string, password string, httpTimeout time.Duration) *SentinelClient {

	return &SentinelClient{
		user:       user,
		password:   password,
		httpClient: &http.Client{},
		searchURL:  "https://scihub.copernicus.eu/dhus/search?q=",
		// searchURL: "https://apihub.copernicus.eu/apihub/search?q=",
		rows:     100,
		dlEngine: sentinel_engine.NewSentinelEngine(user, password, httpTimeout),
	}
}

func (c *SentinelClient) Download(id string, dst string) error {
	if c.dlEngine == nil {
		return fmt.Errorf("no download engine available")
	}

	return c.dlEngine.Download(id, dst)
}

func (c *SentinelClient) IsOnline(id string) (bool, error) {
	if c.dlEngine == nil {
		return false, fmt.Errorf("no download engine available")
	}

	return c.dlEngine.IsOnline(id)
}
