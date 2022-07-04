package service

import (
	"context"
	"errors"
	"fampay/domain"
	videorepository "fampay/repositories/video"
)

var (
	ErrInvalidArgument = errors.New("Invalid Argument")
)

type Service interface {
	Search(ctx context.Context, query string, pageNumber int64, pageItemCount int64) (videos []domain.Video, err error)
}

type service struct {
	videoRepo videorepository.Repository
}

func (svc *service) Search(ctx context.Context, query string, pageNumber int64, pageItemCount int64) (videos []domain.Video, err error) {
	if len(query) < 1 {
		err = ErrInvalidArgument
		return
	}

	if pageNumber < 1 {
		// default pageNumber if not given
		pageNumber = 1
	}

	if pageItemCount < 1 {
		// default pageItemCount if not given
		pageItemCount = 10
	}

	return svc.videoRepo.Search(ctx, query, pageNumber, pageItemCount)
}
