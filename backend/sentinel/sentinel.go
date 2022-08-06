package sentinel_engine

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/cheggaaa/pb/v3"
)

type SentinelEngine struct {
	user       string
	password   string
	httpClient *http.Client
}

// NewSentinelEngine returns a new SentinelEngine
func NewSentinelEngine(user string, password string, httpTimeout time.Duration) SentinelEngine {
	// Function creates new SentinelEngine with given user and password and http client with given timeout.
	// If timeout equals 0 then notimeout is used.

	return SentinelEngine{
		user:     user,
		password: password,
		httpClient: &http.Client{
			Timeout: httpTimeout,
		},
	}
}

func (se SentinelEngine) Download(ProductID string, dst string) error {
	link := fmt.Sprintf("https://scihub.copernicus.eu/dhus/odata/v1/Products('%s')/$value", ProductID)

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return fmt.Errorf("error on create request: %s", err)
	}
	req.SetBasicAuth(se.user, se.password)

	resp, err := se.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error on GET file: %s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == 202 {
		fmt.Printf("Product with product id %s is not ready yet. Triggered offline retrieval.\n", ProductID)
		return fmt.Errorf("file triggered from long-term archive")
	}

	if resp.StatusCode != 200 {
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error on read response body: %s", err)
		}
		return fmt.Errorf("error on GET file: %s", string(bs))
	}

	size, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		panic(err)
	}

	dst_fileName := strings.Trim(strings.TrimSpace(strings.Split(resp.Header.Get("Content-Disposition"), "=")[1]), "\"")

	checkSum := resp.Header.Get("Etag")

	bar := pb.Full.Start64(size)
	bar.Set("prefix", fmt.Sprintf("[ %s ]", dst_fileName))

	barReader := bar.NewProxyReader(resp.Body)
	defer barReader.Close()

	out, err := os.Create(path.Join(dst, dst_fileName))
	if err != nil {
		return fmt.Errorf("error on create local file: %s", err)
	}
	defer out.Close()

	var hashMD5 = md5.New()
	w := io.MultiWriter(out, hashMD5)

	_, err = io.Copy(w, barReader)
	if err != nil {
		return fmt.Errorf("error on saving file: %s", err)
	}
	md5Sum := hashMD5.Sum(nil)

	if checkSum != fmt.Sprintf("%x", md5Sum) {
		return fmt.Errorf("integrity error: checksum mismatch")
	}

	return nil
}
