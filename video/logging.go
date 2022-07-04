package videoservice

import (
	"context"
	"fampay/domain"
	"time"

	"github.com/go-kit/kit/log"
)

type loggingMw struct {
	Service
	logger log.Logger
}

func NewLoggingMw(svc Service, logger log.Logger) *loggingMw {
	return &loggingMw{
		Service: svc,
		logger:  logger,
	}
}

func (l *loggingMw) Search(ctx context.Context, query string, pageNumber int64, pageItemCount int64) (videos []domain.Video, err error) {
	defer func(begin time.Time) {
		l.logger.Log(
			"method", "Search",
			"query", query,
			"pageNumber", pageNumber,
			"pageItemCount", pageItemCount,
			"videos", videos,
			"err", err,
			"took", time.Since(begin))
	}(time.Now())
	return l.Service.Search(ctx, query, pageNumber, pageItemCount)
}
