package main

import (
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
	QUERY       = "vlog"
	RUNNER_FREQ = 10 * time.Second
)

func main() {
	quit := make(chan struct{})
	args := os.Args
	if len(args) < 2 {
		panic("API Key not provide")
	}
	apiKey := args[1]
	svc, err := youtube.NewLb([]string{apiKey})
	if err != nil {
		panic(err)
	}
	initPubAfterDuration := time.Now().Add(-10 * time.Minute)
	runner.StartRunner(RUNNER_FREQ, svc, initPubAfterDuration, QUERY, quit)

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
