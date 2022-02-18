package main

import (
	sentinel "github.com/therox/go-sentinel"
)

func main() {
	client := sentinel.NewClient()

	// Construct OpenAPI Search parameters

	searchParameters := sentinel.OpenAPISearchParams{
		PlatformName: sentinel.Sentinel2,
	}

	client.SearchOData(searchParameters)

}
