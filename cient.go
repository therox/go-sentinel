package sentinel

import (
	"fmt"
	"net/http"
)

type SentinelClient struct {
	Searcher ISentinelSearcher
	dlEngine dlEngine
}

type ISentinelSearcher interface {
	Query(params SearchParameters) (QueryResponse, error)
}

type sentinelSearcher struct {
	user       string
	password   string
	httpClient *http.Client
	searchURL  string
	rows       int
}

func NewSentinelSearcher(user string, password string) ISentinelSearcher {
	return sentinelSearcher{
		user:       user,
		password:   password,
		httpClient: &http.Client{},
		searchURL:  "https://scihub.copernicus.eu/dhus/search?q=",
		// searchURL: "https://apihub.copernicus.eu/apihub/search?q=",
		rows: 100,
	}
}

// func NewClient(user string, password string, httpTimeout time.Duration) *SentinelClient {
func NewClient(searcher ISentinelSearcher, engine dlEngine) *SentinelClient {

	sc := &SentinelClient{
		Searcher: searcher,
		dlEngine: engine,
	}
	return sc
}

func (c *SentinelClient) Download(id string, dst string) (string, error) {
	if c.dlEngine == nil {
		return "", fmt.Errorf("no download engine available")
	}

	return c.dlEngine.Download(id, dst)
}

func (c *SentinelClient) IsOnline(id string) (bool, error) {
	if c.dlEngine == nil {
		return false, fmt.Errorf("no download engine available")
	}

	return c.dlEngine.IsOnline(id)
}
