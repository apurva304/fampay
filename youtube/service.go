package youtube

import (
	"fmt"
	"net/http"
	"time"

	"google.golang.org/api/googleapi/transport"
	youtube "google.golang.org/api/youtube/v3"
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
	req := svc.ytClient.Search.List([]string{"snippet"}).
		Q(query).
		MaxResults(25).
		Type("video").
		PublishedAfter(publishedAfter.Format(time.RFC3339)).
		Order("date")

	res, err := req.Do()
	if err != nil {
		return err
	}
	fmt.Println(len(res.Items))
	if len(res.Items) < 1 {
		fmt.Println("NOT FOUND")
		return nil
	}

	return
}
