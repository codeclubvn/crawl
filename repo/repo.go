package repo

import (
	"crawl/model"
	"go.mongodb.org/mongo-driver/bson"
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
	CreateBulk(input []model.CompanyInfo) (err error)
	GetOne(filter bson.M) (model.CompanyInfo, error)
	GetList(filter bson.M, limit int64, skip int64) ([]model.CompanyInfo, error)
}
