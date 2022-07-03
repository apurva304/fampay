package main

import (
	"fampay/youtube"
	"os"
	"time"
)

var ()

func main() {
	args := os.Args
	if len(args) < 2 {
		panic("API Key not provide")
	}
	apiKey := args[1]
	svc, err := youtube.NewService(apiKey)
	if err != nil {
		panic(err)
	}

	err = svc.Search("music", time.Now().Add(-60*time.Minute))
	if err != nil {
		panic(err)
	}
}
