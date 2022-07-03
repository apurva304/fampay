package youtube

import (
	"time"

	"github.com/go-kit/kit/log"
)

type loggingMW struct {
	Service
	logger log.Logger
}

func (l *loggingMW) Search(query string, publishedAfter time.Time) (err error) {
	defer func(begin time.Time) {
		l.logger.Log(
			"method", "Search",
			"query", query,
			"publishAfter", publishedAfter,
			"err", err,
			"took", time.Since(begin))
	}(time.Now())
	return l.Service.Search(query, publishedAfter)
}
