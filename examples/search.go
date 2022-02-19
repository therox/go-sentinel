package main

import (
	"log"

	sentinel "github.com/therox/go-sentinel"
)

func main() {
	client := sentinel.NewClient()

	// Construct OpenAPI Search parameters

	searchParameters := sentinel.SearchParameters{
		Platforms: []sentinel.Platform{sentinel.PlanformSentinel2},
	}

	_, err := client.Query(searchParameters)
	if err != nil {
		log.Fatal(err)
	}

}
