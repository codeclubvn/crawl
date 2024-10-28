package repo

import (
	"context"
	"crawl/model"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *Repo) CreateEmailIndex() error {
	collection := r.db.Collection("companies")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Tạo unique index cho email field
	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("error creating unique index: %v", err)
	}

	return nil
}

func (r *Repo) CreateBulk(input []model.CompanyInfo) error {
	if len(input) == 0 {
		return nil
	}

	if err := r.CreateEmailIndex(); err != nil {
		return err
	}
	collection := r.db.Collection("companies")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Tạo bulk write operations
	var operations []mongo.WriteModel
	for _, company := range input {
		if company.Email == "" {
			continue
		}
		op := mongo.NewInsertOneModel().SetDocument(company)
		operations = append(operations, op)
	}

	// Thực hiện bulk write với ordered=false để bỏ qua các lỗi duplicate
	opts := options.BulkWrite().SetOrdered(false)
	_, err := collection.BulkWrite(ctx, operations, opts)
	if err != nil {
		// Kiểm tra nếu là lỗi duplicate thì bỏ qua
		if !mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("error bulk writing companies: %v", err)
		}
	}

	return nil
}

func (r *Repo) GetOne(filter bson.M) (model.CompanyInfo, error) {
	collection := r.db.Collection("companies")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var company model.CompanyInfo
	err := collection.FindOne(ctx, filter).Decode(&company)
	if err != nil {
		return company, fmt.Errorf("error finding company: %v", err)
	}

	return company, nil
}

func (r *Repo) GetList(filter bson.M, limit int64, skip int64) ([]model.CompanyInfo, error) {
	collection := r.db.Collection("companies")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Tạo options để limit và skip
	opts := options.Find().SetLimit(limit).SetSkip(skip)
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("error finding companies: %v", err)
	}
	defer cursor.Close(ctx)

	var companies []model.CompanyInfo
	if err := cursor.All(ctx, &companies); err != nil {
		return nil, fmt.Errorf("error decoding companies: %v", err)
	}

	return companies, nil
}
