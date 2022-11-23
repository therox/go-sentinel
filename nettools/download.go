package tools

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadFile(fromURL string, toFile string) error {

	out, err := os.Create(toFile)
	if err != nil {
		return fmt.Errorf("error on create local index file: %s", err)
	}
	defer out.Close()

	resp, err := http.Get(fromURL)
	if err != nil {
		return fmt.Errorf("error on get index file: %s", err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("error on saving index file: %s", err)
	}
	return nil
}
