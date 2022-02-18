package main

import (
	sentinel "github.com/therox/go-sentinel"
)

func main() {
	client := sentinel.NewClient()

	// Construct OpenAPI Search parameters

	searchParameters := sentinel.SearchParameters{
		Platforms: []sentinel.Platform{sentinel.PlanformSentinel2},
	}

	client.Query(searchParameters)

}
