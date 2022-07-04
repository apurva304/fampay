package mongovideorepository

import (
	"context"
	"errors"
	"fampay/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrVideoAlreadyExists = errors.New("Video Already Exists")
)

const (
	COLLECTION = "video"
)

type repository struct {
	client *mongo.Client
	dbName string
}

func New(client *mongo.Client, dbName string) *repository {
	index := mongo.IndexModel{
		Keys: bson.D{
			{"title", "text"},
			{"description", "text"},
		},
	}
	_, err := client.Database(dbName).Collection(COLLECTION).Indexes().CreateOne(context.TODO(), index)
	if err != nil {
		// mongo.
	}
	return &repository{}
}

func (repo *repository) Add(ctx context.Context, video domain.Video) (err error) {
	_, err = repo.client.Database(repo.dbName).Collection(COLLECTION).InsertOne(ctx, video)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			err = ErrVideoAlreadyExists
		}
		return
	}
	return
}
func (repo *repository) AddBulk(ctx context.Context, vidoes []domain.Video) (err error) {
	var data []interface{}
	for _, v := range vidoes {
		data = append(data, v)
	}

	_, err = repo.client.Database(repo.dbName).Collection(COLLECTION).InsertMany(ctx, data)
	if err != nil {
		return
	}
	return
}
func (repo *repository) Search(ctx context.Context, query string, pageNumber int64, pageItemCount int64) (videos []domain.Video, err error) {
	q := bson.M{
		"$text": bson.M{
			"$search": query,
		},
	}

	skip := (pageNumber - 1) * pageItemCount
	opts := options.Find().SetLimit(pageItemCount)

	if skip > 0 {
		opts = opts.SetSkip(skip)
	}

	curr, err := repo.client.Database(repo.dbName).Collection(COLLECTION).Find(ctx, q)
	if err != nil {
		return
	}

	err = curr.All(ctx, &videos)
	if err != nil {
		return
	}

	return
}
