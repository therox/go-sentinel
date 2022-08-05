package sentinel

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
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

func (c *SentinelClient) Download(id string, dst string) error {
	fmt.Println("GOT ID: ", id)
	link := fmt.Sprintf("https://scihub.copernicus.eu/dhus/odata/v1/Products('%s')/$value", id)

	fmt.Printf("Downloading file %s to %s\n", link, dst)

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return fmt.Errorf("error on create request: %s", err)
	}
	req.SetBasicAuth(c.user, c.password)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error on get file: %s", err)
	}

	fmt.Println("Headers: ", resp.Header.Get("Content-Disposition"))
	dst_fileName := strings.Trim(strings.TrimSpace(strings.Split(resp.Header.Get("Content-Disposition"), "=")[1]), "\"")
	defer resp.Body.Close()

	out, err := os.Create(path.Join(dst, dst_fileName))
	if err != nil {
		return fmt.Errorf("error on create local file: %s", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("error on saving file: %s", err)
	}
	fmt.Println("Download complete")
	return nil
}
