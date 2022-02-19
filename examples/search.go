package main

import (
	"fmt"
	"log"

	sentinel "github.com/therox/go-sentinel"
)

func main() {
	client := sentinel.NewClient()

	// Construct OpenAPI Search parameters

	searchParameters := sentinel.SearchParameters{
		Platforms: []sentinel.Platform{sentinel.PlanformSentinel2},
		TileIDs:   []string{"37UCU", "37UCT"},
	}

	res, err := client.Query(searchParameters)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(res.Feed.Entries))

}
