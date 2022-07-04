package main

import (
	"fampay/config"
	runner "fampay/jobrunner"
	"fampay/youtube"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oklog/oklog/pkg/group"
)

const (
	RUNNER_FREQ = 10 * time.Second
	CONFIG_FILE = "config.json"
)

func main() {
	file, err := os.Open(CONFIG_FILE)
	if err != nil {
		panic(err)
	}

	conf, err := config.New(file)
	if err != nil {
		panic(err)
	}

	quit := make(chan struct{})
	svc, err := youtube.NewLb(conf.ApiKey)
	if err != nil {
		panic(err)
	}
	initPubAfterDuration := time.Now().Add(-10 * time.Minute)
	runner.StartRunner(RUNNER_FREQ, svc, initPubAfterDuration, conf.Query, nil, quit)

	var g group.Group
	// Interupt
	{
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
			select {
			case sig := <-c:
				fmt.Println("serverstate", fmt.Sprintf("received signal %s", sig))
				quit <- struct{}{} // shutdown the runner gorountine
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}
	fmt.Printf("exiting... error: %v\n", g.Run())
}
