package sentinel

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	pb "github.com/cheggaaa/pb/v3"
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
		log.Fatalf("Please provide SENTINEL_CREDENTIALS!")
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

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return fmt.Errorf("error on create request: %s", err)
	}
	req.SetBasicAuth(c.user, c.password)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error on get file: %s", err)
	}

	size, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		panic(err)
	}
	dst_fileName := strings.Trim(strings.TrimSpace(strings.Split(resp.Header.Get("Content-Disposition"), "=")[1]), "\"")
	defer resp.Body.Close()

	bar := pb.Full.Start64(size)
	barReader := bar.NewProxyReader(resp.Body)

	out, err := os.Create(path.Join(dst, dst_fileName))
	if err != nil {
		return fmt.Errorf("error on create local file: %s", err)
	}
	defer out.Close()

	_, err = io.Copy(out, barReader)
	if err != nil {
		return fmt.Errorf("error on saving file: %s", err)
	}
	bar.Finish()
	return nil
}
