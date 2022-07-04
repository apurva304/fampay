package videorepository

import (
	"context"
	"fampay/domain"
)

type Repository interface {
	Add(ctx context.Context, video domain.Video) (err error)
	AddBulk(ctx context.Context, vidoes []domain.Video) (err error)
	Search(ctx context.Context, query string, pageNumber int64, pageItemCount int64) (videos []domain.Video, err error)
}
