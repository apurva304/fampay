package jobrunner

import (
	"fampay/youtube"
	"log"
	"time"
)

type runner struct {
	ticker               *time.Ticker
	svc                  youtube.Service
	lastSuccessFetchTime time.Time
	query                string
	quit                 chan struct{}
}

func StartRunner(runDuration time.Duration, svc youtube.Service, publishAfter time.Time, query string, quit chan struct{}) {
	r := &runner{
		ticker:               time.NewTicker(runDuration),
		svc:                  svc,
		lastSuccessFetchTime: publishAfter,
		quit:                 quit,
	}
	go r.run()
}

func (r *runner) run() {
	for {
		select {
		case <-r.ticker.C:
			err := r.svc.Search(r.query, r.lastSuccessFetchTime)
			switch err {
			case nil:
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
