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
	Search(ctx context.Context, query string) (videos []domain.Video, err error)
}

type service struct {
	videoRepo videorepository.Repository
}

func (svc *service) Search(ctx context.Context, query string) (videos []domain.Video, err error) {
	if len(query) < 1 {
		err = ErrInvalidArgument
		return
	}

	return svc.videoRepo.Search(query)
}
