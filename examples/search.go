package main

import (
	sentinel "github.com/therox/go-sentinel"
)

func main() {
	client := sentinel.NewClient()
	client.SearchOData()
}
