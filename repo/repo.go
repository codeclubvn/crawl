package repo

import (
	"crawl/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repo struct {
	db *mongo.Database
}

func NewRepo(db *mongo.Database) *Repo {
	return &Repo{
		db: db,
	}
}

type IRepo interface {
	GetUserByName(name string) (model.User, error)
}
