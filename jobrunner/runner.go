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
}

func StartRunner(runDuration time.Duration, svc youtube.Service, publishAfter time.Time, query string) {
	r := &runner{
		ticker:               time.NewTicker(runDuration),
		svc:                  svc,
		lastSuccessFetchTime: publishAfter,
	}
	r.run()
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
				break
			default:
				log.Println("Error While Fetching youtube data", err)
				continue
			}
		}
	}
}
