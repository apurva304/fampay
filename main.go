package main

import (
	"context"
	"fampay/config"
	runner "fampay/jobrunner"
	videorepository "fampay/repositories/video"
	mongovideorepo "fampay/repositories/video/mongo"
	videoservice "fampay/video"
	"fampay/youtube"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/oklog/oklog/pkg/group"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	logger := log.NewJSONLogger(os.Stdout)

	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(conf.MongoUri))
	if err != nil {
		panic(err)
	}

	var videoRepo videorepository.Repository
	{
		videoRepo, err = mongovideorepo.New(mongoClient, conf.DBName)
		if err != nil {
			panic(err)
		}
		videoRepo = videorepository.NewLoggingMw(videoRepo, log.With(logger, "repository", "videorepsoitry"))
	}

	quit := make(chan struct{})

	var youtubeSvc youtube.Service
	{
		youtubeSvc, err = youtube.NewLb(conf.ApiKey)
		if err != nil {
			panic(err)
		}
		youtubeSvc = youtube.NewLoggingMw(youtubeSvc, log.With(logger, "service", "youtube"))
	}

	initPubAfterDuration := time.Now().Add(-10 * time.Minute)
	runner.StartRunner(RUNNER_FREQ, youtubeSvc, initPubAfterDuration, conf.Query, videoRepo, quit, log.With(logger, "service", "jobrunner"))

	var svc videoservice.Service
	{
		svc = videoservice.New(videoRepo)
		videoservice.NewLoggingMw(svc, log.With(logger, "service", "video"))
	}

	httpRouter := videoservice.MakeHandler(svc, logger)

	var g group.Group
	// Game HTTP server
	{
		httpServer := &http.Server{
			Addr:    ":" + strconv.Itoa(conf.HttpPort),
			Handler: httpRouter,
		}
		g.Add(func() error {
			logger.Log("serverstate", "starting game http server. port:%d \n", conf.HttpPort)
			return httpServer.ListenAndServe()
		}, func(error) {
			ctx, cancle := context.WithTimeout(context.Background(), time.Second*5)
			defer cancle()
			httpServer.Shutdown(ctx)
		})
	}
	// Interupt
	{
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
			select {
			case sig := <-c:
				fmt.Println("serverstate", fmt.Sprintf("received signal %s", sig))
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}
	fmt.Printf("exiting... error: %v\n", g.Run())
	quit <- struct{}{} // shutdown the runner gorountine
}
