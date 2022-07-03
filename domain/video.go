package domain

type Video struct {
	Id            string
	Title         string
	Description   string
	ChannelId     string
	ChannelName   string
	PublishedAt   string
	ThumbnailUrls []string
}
