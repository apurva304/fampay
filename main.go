package main

import (
	"fampay/youtube"
	"os"
)

var ()

func main() {
	args := os.Args
	if len(args) < 2 {
		panic("API Key not provide")
	}
	apiKey := args[1]
	svc, err := youtube.NewLb([]string{apiKey})
	if err != nil {
		panic(err)
	}

}
