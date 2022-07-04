package jobrunner

import (
	videorepository "fampay/repositories/video"
	"fampay/youtube"
	"log"
	"time"
)

type runner struct {
	ticker               *time.Ticker
	svc                  youtube.Service
	lastSuccessFetchTime time.Time
	query                string
	videoRepo            videorepository.Repository
	quit                 chan struct{}
}

func StartRunner(runDuration time.Duration, svc youtube.Service, publishAfter time.Time, query string, videoRepo videorepository.Repository, quit chan struct{}) {
	r := &runner{
		ticker:               time.NewTicker(runDuration),
		svc:                  svc,
		lastSuccessFetchTime: publishAfter,
		videoRepo:            videoRepo,
		quit:                 quit,
	}
	go r.run()
}

func (r *runner) run() {
	for {
		select {
		case <-r.ticker.C:
			videos, err := r.svc.Search(r.query, r.lastSuccessFetchTime)
			switch err {
			case nil:
				err = r.videoRepo.AddBulk(videos)
				if err != nil {
					log.Println("Error While Adding Data", err)
					continue
				}
				r.lastSuccessFetchTime = time.Now()
			case youtube.ErrNotFound:
			case youtube.ErrQuotaExceeded:
				//limit exceeded for all the provided keys
				log.Fatal("Quota Limit Exceeded for all the provided keys")
			default:
				log.Println("Error While Fetching youtube data", err)
				continue
			}
		case <-r.quit:
			return
		}
	}
}
