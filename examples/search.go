package main

import (
	"fmt"
	"log"

	sentinel "github.com/therox/go-sentinel"
)

func main() {
	client := sentinel.NewClient()

	// Construct OpenAPI Search parameters
	tiles := []string{"36UYA", "36UYB", "36UYC", "36UYD", "36UYE", "37TDK", "37TDL", "37TEL", "37TEM", "37UCR", "37UCS", "37UCT", "37UCU", "37UCV", "37UDR", "37UDS", "37UDT", "37UDU", "37UDV", "37UEU", "37UEV", "37UFT", "37UFU", "37UFV", "37UGT", "37UGU", "38ULD", "44UPG", "45UUB"}

	resCount := 0
	entries := make([]sentinel.QueryEntryResponse, 0)
	for _, tile := range tiles {
		searchParameters := sentinel.SearchParameters{
			Platforms: []sentinel.Platform{sentinel.PlanformSentinel2},
			// Filenames:     []string{fmt.Sprintf("*%s*", tile)},
			TileIDs: []string{tile},
			// ProductTypes: []string{"S2MSI2A", "S2MS2Ap"},
			// ProductTypes:  []string{"S2MSI1C"},
			BeginPosition: "[2022-01-01T00:00:00.000Z TO NOW]",
		}

		res, err := client.Query(searchParameters)
		if err != nil {
			log.Fatal(err)
		}
		resCount += res.Feed.TotalResults
		entries = append(entries, res.Feed.Entries...)
	}
	fmt.Printf("Total found %d items\n", resCount)

	if resCount > 0 {
		// fmt.Printf("%+v\n", entries[0])
		// fmt.Printf("%+v\n", string(entries[0].Date))
		client.Download(entries[0].ID, "/tmp")
		// fmt.Printf("%+v\n", entries[0].Link[0])
		// tools.DownloadFile(entries[0].Link[0].HREF, "/tmp/1.zip")

	}

}
