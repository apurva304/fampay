package youtube

import (
	"errors"
	"net/http"
	"time"

	"fampay/domain"

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
	MAX_RESULT      = 100
)

type Service interface {
	Search(query string, publishedAfter time.Time) (videos []domain.Video, err error)
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

func (svc *service) Search(query string, publishedAfter time.Time) (videos []domain.Video, err error) {
	req := svc.ytClient.Search.List([]string{LIST_QUERY_PART}).
		Q(query).
		MaxResults(MAX_RESULT).
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
		err = ErrNotFound
		return
	}

	for _, item := range res.Items {
		video, valid := getVideo(item)
		if !valid {
			continue
		}
		videos = append(videos, video)
	}
	return
}

func getVideo(item *youtube.SearchResult) (video domain.Video, valid bool) {
	if item == nil {
		return video, false
	}
	if item.Id != nil {
		video.Id = item.Id.VideoId
	}
	snippet := item.Snippet

	if snippet != nil {
		valid = true
		video.Title = snippet.Title
		video.Description = snippet.Description
		video.ChannelId = snippet.ChannelId
		video.ChannelName = snippet.ChannelTitle
		video.PublishedAt = snippet.PublishedAt
	}

	thumbnail := snippet.Thumbnails

	if thumbnail == nil {
		return video, true
	}

	if thumbnail.Standard != nil && len(thumbnail.Standard.Url) > 0 {
		video.ThumbnailUrls = append(video.ThumbnailUrls, thumbnail.Standard.Url)
	}
	if thumbnail.High != nil && len(thumbnail.High.Url) > 0 {
		video.ThumbnailUrls = append(video.ThumbnailUrls, thumbnail.High.Url)
	}
	if thumbnail.Maxres != nil && len(thumbnail.Maxres.Url) > 0 {
		video.ThumbnailUrls = append(video.ThumbnailUrls, thumbnail.Maxres.Url)
	}
	if thumbnail.Medium != nil && len(thumbnail.Medium.Url) > 0 {
		video.ThumbnailUrls = append(video.ThumbnailUrls, thumbnail.Medium.Url)
	}
	if thumbnail.Default != nil && len(thumbnail.Default.Url) > 0 {
		video.ThumbnailUrls = append(video.ThumbnailUrls, thumbnail.Default.Url)
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
