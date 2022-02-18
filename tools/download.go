package tools

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadFile(fromURL string, toFile string) error {
	fmt.Printf("Downloading file %s to %s\n", fromURL, toFile)
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
	fmt.Println("Download ready")
	return nil
}

// gzReader, err := gzip.NewReader(resp.Body)
// if err != nil {
// 	log.Fatalf("Error on Instatiate GZ reader from resp: %s", err)
// }
// defer gzReader.Close()
