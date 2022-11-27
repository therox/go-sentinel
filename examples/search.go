package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	sentinel "github.com/therox/go-sentinel"
	sentinel_engine "github.com/therox/go-sentinel/backend/sentinel"
)

func main() {

	credentials := strings.Split(os.Getenv("SENTINEL_CREDENTIALS"), ":")
	if len(credentials) < 2 {
		log.Fatalf("Please provide SENTINEL_CREDENTIALS!")
	}

	searcher := sentinel.NewSentinelSearcher(credentials[0], credentials[1])

	dlEngine := sentinel_engine.NewSentinelEngine(credentials[0], credentials[1], 60*time.Minute)

	client := sentinel.NewClient(searcher, dlEngine)

	// Construct OpenAPI Search parameters
	// tiles := []string{"36UYA", "36UYB", "36UYC", "36UYD", "36UYE", "37TDK", "37TDL", "37TEL", "37TEM", "37UCR", "37UCS", "37UCT", "37UCU", "37UCV", "37UDR", "37UDS", "37UDT", "37UDU", "37UDV", "37UEU", "37UEV", "37UFT", "37UFU", "37UFV", "37UGT", "37UGU", "38ULD", "44UPG", "45UUB"}
	tiles := []string{"36UYA"}

	resCount := 0
	entries := make([]sentinel.QueryEntryResponse, 0)
	for _, tile := range tiles {
		print("Searching for tile " + tile + "...")
		et := time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC)
		searchParameters := sentinel.SearchParameters{
			Platforms: []sentinel.Platform{sentinel.PlanformSentinel2},
			// Filenames:     []string{fmt.Sprintf("*%s*", tile)},
			// TileIDs:      []string{tile},
			Footprint:    "POLYGON((35.8207689450001 50.518113703,37.3657427630001 50.4703606300001,37.2773703820001 49.485620291,35.763545148 49.531744406,35.8207689450001 50.518113703))",
			ProductTypes: []string{"S2MSI2A", "S2MS2Ap"},
			BeginDate:    time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:      &et,
		}

		res, err := client.Searcher.Query(searchParameters)
		if err != nil {
			log.Fatal(err)
		}
		resCount += res.Feed.TotalResults
		entries = append(entries, res.Feed.Entries...)
	}
	fmt.Printf("Total found %d items\n", len(entries))

	if resCount > 0 {
		for _, entry := range entries {
			isOnline, err := client.IsOnline(entry.ID)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("[%s] %t:%s\n", entry.FileName, isOnline, entry.BeginPosition)
			if isOnline {
				fName, err := client.Download(entry.GetID(), "/tmp")
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("Downloaded to " + fName)
			}
		}
	}
}
