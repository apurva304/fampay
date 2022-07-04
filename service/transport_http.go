package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/go-kit/kit/transport"
	"github.com/gorilla/mux"
)

var (
	ErrBadRequest = errors.New("Bad Request")
)

func MakeHandler(svc Service, logger kitlog.Logger) {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	getVideoHandler := kithttp.NewServer(
		makeGetVideoEndpoint(svc),
		decodeGetVideoRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()
	r.Handle("/videos", getVideoHandler).Methods(http.MethodGet)
}

func decodeGetVideoRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req getVideoRequest
	queryMap := r.URL.Query()
	q, ok := queryMap["query"]
	if !ok {
		err = ErrBadRequest
		return
	}
	if len(q) < 1 {
		err = ErrBadRequest
		return
	}
	req.Query = q[0]

	return req, nil
}

type errorer interface {
	error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
