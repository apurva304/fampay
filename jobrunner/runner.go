package jobrunner

import (
	"context"
	videorepository "fampay/repositories/video"
	"fampay/youtube"

	"time"

	"github.com/go-kit/kit/log"
)

type runner struct {
	ticker               *time.Ticker
	svc                  youtube.Service
	lastSuccessFetchTime time.Time
	query                string
	videoRepo            videorepository.Repository
	logger               log.Logger
	quit                 chan struct{}
}

func StartRunner(runDuration time.Duration, svc youtube.Service, publishAfter time.Time, query string, videoRepo videorepository.Repository, quit chan struct{}, logger log.Logger) {
	r := &runner{
		ticker:               time.NewTicker(runDuration),
		svc:                  svc,
		lastSuccessFetchTime: publishAfter,
		query:                query,
		videoRepo:            videoRepo,
		logger:               logger,
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
				err = r.videoRepo.AddBulk(context.TODO(), videos)
				if err != nil {
					r.logger.Log("Error While Adding Data", err)
					continue
				}
				r.lastSuccessFetchTime = time.Now()
			case youtube.ErrNotFound:
			case youtube.ErrQuotaExceeded:
				//limit exceeded for all the provided keys
				r.logger.Log("Quota Limit Exceeded for all the provided keys")
				continue
			default:
				r.logger.Log("Error While Fetching youtube data", err)
				continue
			}
		case <-r.quit:
			return
		}
	}
}
