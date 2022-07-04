package videoservice

import (
	"context"
	"fampay/domain"

	"github.com/go-kit/kit/endpoint"
)

type searchRequest struct {
	Query         string `json:"query"`
	PageNumber    int64  `json:"pageNumber"`
	PageItemCount int64  `json:"pageItemCount"`
}

type searchResponse struct {
	Videos []domain.Video
	Err    error
}

func (res searchResponse) error() string {
	return res.Err.Error()
}

func makeSearchEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(searchRequest)
		videos, err := svc.Search(ctx, req.Query, req.PageNumber, req.PageItemCount)
		return searchResponse{Videos: videos, Err: err}, nil
	}
}
