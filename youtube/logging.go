package youtube

import (
	"fampay/domain"
	"time"

	"github.com/go-kit/kit/log"
)

type loggingMW struct {
	Service
	logger log.Logger
}

func NewLoggingMw(svc Service, logger log.Logger) *loggingMW {
	return &loggingMW{
		Service: svc,
		logger:  logger,
	}
}

func (l *loggingMW) Search(query string, publishedAfter time.Time) (videos []domain.Video, err error) {
	defer func(begin time.Time) {
		l.logger.Log(
			"method", "Search",
			"query", query,
			"publishAfter", publishedAfter,
			"len(videos)", len(videos),
			"err", err,
			"took", time.Since(begin))
	}(time.Now())
	return l.Service.Search(query, publishedAfter)
}
