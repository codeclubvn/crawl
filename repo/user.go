package repo

import (
	"context"
	"crawl/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *Repo) GetUserByName(name string) (model.User, error) {
	// Tìm kiếm bản ghi trong bảng "users"từ request
	table := r.db.Collection("users")
	var user model.User
	query := bson.M{"username": name}
	if err := table.FindOne(context.Background(), query).Decode(&user); err != nil {
		return model.User{}, err
	}

	return user, nil
}
