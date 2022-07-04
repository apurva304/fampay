package videorepository

import "fampay/domain"

type Repository interface {
	Add(video domain.Video) (err error)
	AddBulk(vidoes []domain.Video) (err error)
	Search(query string) (videos []domain.Video, err error)
}
