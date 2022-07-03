package youtube

import (
	"errors"
	"net/http"
	"time"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/googleapi/transport"
	youtube "google.golang.org/api/youtube/v3"
)

var (
	ErrNotFound      = errors.New("Videos Not Found")
	ErrQuotaExceeded = errors.New("Quota Exceeded")
)

const (
	QOUTA_EXCEEDED  = "quotaExceeded"
	LIST_QUERY_PART = "snippet"
	TYPE            = "video"
	ORDER           = "date"
)

type Service interface {
	Search(query string, publishedAfter time.Time) (err error)
}

type service struct {
	ytClient *youtube.Service
}

func NewService(apiKey string) (*service, error) {
	client := &http.Client{
		Transport: &transport.APIKey{Key: apiKey},
	}

	ytClient, err := youtube.New(client)
	if err != nil {
		return nil, err
	}

	svc := &service{
		ytClient: ytClient,
	}

	return svc, nil
}

func (svc *service) Search(query string, publishedAfter time.Time) (err error) {
	req := svc.ytClient.Search.List([]string{LIST_QUERY_PART}).
		Q(query).
		MaxResults(25).
		Type(TYPE).
		PublishedAfter(publishedAfter.Format(time.RFC3339)).
		Order(ORDER)

	res, err := req.Do()
	switch {
	case err == nil:
		// continue below
	case checkQuotaExceeded(err):
		err = ErrQuotaExceeded
		return
	default:
		return
	}
	if len(res.Items) < 1 {
		return ErrNotFound
	}

	return
}

func checkQuotaExceeded(err error) (ok bool) {
	var gApiErr *googleapi.Error
	if errors.As(err, &gApiErr) {
		if len(gApiErr.Errors) > 0 && gApiErr.Errors[0].Reason == QOUTA_EXCEEDED {
			ok = true
		}
	}
	return
}
