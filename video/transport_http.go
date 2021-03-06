package videoservice

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/go-kit/kit/transport"
	"github.com/gorilla/mux"
)

var (
	ErrBadRequest = errors.New("Bad Request")
)

func MakeHandler(svc Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	getVideoHandler := kithttp.NewServer(
		makeSearchEndpoint(svc),
		decodeSearchRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()
	r.Handle("/videos", getVideoHandler).Methods(http.MethodGet)

	return r
}

func decodeSearchRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req searchRequest
	queryMap := r.URL.Query()
	queryValue, ok := queryMap["query"]
	if !ok {
		err = ErrBadRequest
		return
	}
	if len(queryValue) < 1 {
		err = ErrBadRequest
		return
	}
	req.Query = queryValue[0]

	// pageNumber is not mandatory
	pageNumber, _ := queryMap["pageNumber"]
	if len(pageNumber) > 0 {
		req.PageNumber, err = strconv.ParseInt(pageNumber[0], 10, 64)
		if err != nil {
			err = ErrBadRequest
			return
		}
	}
	// pageItemCount is not mandatory
	pageItemCount, _ := queryMap["pageItemCount"]
	if len(pageItemCount) > 0 {
		req.PageItemCount, err = strconv.ParseInt(pageItemCount[0], 10, 64)
		if err != nil {
			err = ErrBadRequest
			return
		}
	}

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
