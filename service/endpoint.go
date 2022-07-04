package service

import (
	"context"
	"fampay/domain"

	"github.com/go-kit/kit/endpoint"
)

type getVideoRequest struct {
	Query string `json:"query"`
}

type getVideoResponse struct {
	Videos []domain.Video
	Err    error
}

func (res getVideoResponse) error() string {
	return res.Err.Error()
}

func makeGetVideoEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getVideoRequest)
		videos, err := svc.Search(ctx, req.Query)
		return getVideoResponse{Videos: videos, Err: err}, nil
	}
}
