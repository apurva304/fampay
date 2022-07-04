package videorepository

import (
	"context"
	"fampay/domain"
	"time"

	"github.com/go-kit/kit/log"
)

type loggingMw struct {
	Repository
	logger log.Logger
}

func NewLoggingMw(repo Repository, logger log.Logger) *loggingMw {
	return &loggingMw{
		Repository: repo,
		logger:     logger,
	}
}

func (repo *loggingMw) Add(ctx context.Context, video domain.Video) (err error) {
	defer func(begin time.Time) {
		repo.logger.Log(
			"method", "Add",
			"video", video,
			"err", err,
			"took", time.Since(begin))
	}(time.Now())
	return repo.Repository.Add(ctx, video)
}
func (repo *loggingMw) AddBulk(ctx context.Context, vidoes []domain.Video) (err error) {
	defer func(begin time.Time) {
		repo.logger.Log(
			"method", "AddBulk",
			"len(videos)", len(vidoes),
			"err", err,
			"took", time.Since(begin))
	}(time.Now())
	return repo.Repository.AddBulk(ctx, vidoes)
}
func (repo *loggingMw) Search(ctx context.Context, query string, pageNumber int64, pageItemCount int64) (videos []domain.Video, err error) {
	defer func(begin time.Time) {
		repo.logger.Log(
			"method", "Search",
			"query", query,
			"pageNumber", pageNumber,
			"pageItemCount", pageItemCount,
			"err", err,
			"took", time.Since(begin))
	}(time.Now())
	return repo.Repository.Search(ctx, query, pageNumber, pageItemCount)
}
