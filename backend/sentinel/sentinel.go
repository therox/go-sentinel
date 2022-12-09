package sentinel_engine

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type (
	ErrFileTriggered struct {
		productID string
	}
	ErrIntegrityError struct {
		productID string
	}
)

func (e ErrFileTriggered) Error() string {
	return fmt.Sprintf("file triggered from long-term archive: %s", e.productID)
}

func (e ErrIntegrityError) Error() string {
	return fmt.Sprintf("dataset %s integrity error: checksum mismatch", e.productID)
}

type SentinelEngine struct {
	user       string
	password   string
	httpClient *http.Client
	dhusURL    string
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
		dhusURL: "https://scihub.copernicus.eu/dhus/odata/v1",
	}
}

func (se SentinelEngine) getURL(product_id string, suffix string) string {
	return fmt.Sprintf("%s/Products('%s')/%s", se.dhusURL, product_id, suffix)
}

func (se SentinelEngine) Download(productID string, dst string) (string, error) {
	filePath := ""
	link := se.getURL(productID, "$value")

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return filePath, fmt.Errorf("error on create request: %s", err)
	}
	req.SetBasicAuth(se.user, se.password)

	resp, err := se.httpClient.Do(req)
	if err != nil {
		return filePath, fmt.Errorf("error on GET file: %s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == 202 {
		return filePath, ErrFileTriggered{productID: productID}
	}

	if resp.StatusCode != 200 {

		return filePath, fmt.Errorf("%d:%s", resp.StatusCode, resp.Header.Get("Cause-Message"))
	}

	_, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return filePath, fmt.Errorf("error on parse Content-Length: %s", err)
	}

	dst_fileName := strings.Trim(strings.TrimSpace(strings.Split(resp.Header.Get("Content-Disposition"), "=")[1]), "\"")

	checkSum := resp.Header.Get("Etag")

	filePath = path.Join(dst, dst_fileName)
	out, err := os.Create(filePath)
	if err != nil {
		return filePath, fmt.Errorf("error on create local file: %s", err)
	}
	defer out.Close()

	hashMD5 := md5.New()
	w := io.MultiWriter(out, hashMD5)

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return filePath, fmt.Errorf("error on saving file: %s", err)
	}

	if checkSum != fmt.Sprintf("%x", hashMD5.Sum(nil)) {
		os.RemoveAll(filePath)
		return filePath, fmt.Errorf(ErrIntegrityError{productID: productID}.Error())
	}
	fmt.Println("File downloaded to ", filePath)

	return filePath, nil
}

func (se SentinelEngine) IsOnline(productID string) (bool, error) {
	link := se.getURL(productID, "Online/$value")

	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return false, fmt.Errorf("error on create request: %s", err)
	}
	req.SetBasicAuth(se.user, se.password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("error on GET Online status: %s", err)
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("error on read response body: %s", err)
	}

	return string(bs) == "true", nil

}
