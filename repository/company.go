package repository

import (
	"context"
	"crawl/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ICompanyRepository interface {
	CreateOne(ctx context.Context, company *domain.Company) error
	DeleteOne(ctx context.Context, id primitive.ObjectID) error
	UpdateOne(ctx context.Context, id primitive.ObjectID, company *domain.Company) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Company, error)
	GetByName(ctx context.Context, name string) (*domain.Company, error)
	GetAll(ctx context.Context) ([]domain.Company, error)
}

type companyRepository struct {
	companyCollection string
	db                *mongo.Database
}

func NewCompanyRepository(companyCollection string, db *mongo.Database) ICompanyRepository {
	return &companyRepository{companyCollection: companyCollection, db: db}
}

func (c *companyRepository) CreateOne(ctx context.Context, company *domain.Company) error {
	companyCollection := c.db.Collection(c.companyCollection)

	_, err := companyCollection.InsertOne(ctx, company)
	if err != nil {
		return err
	}

	return nil
}

func (c *companyRepository) DeleteOne(ctx context.Context, id primitive.ObjectID) error {
	companyCollection := c.db.Collection(c.companyCollection)

	filter := bson.M{"_id": id}
	_, err := companyCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (c *companyRepository) UpdateOne(ctx context.Context, id primitive.ObjectID, company *domain.Company) error {
	companyCollection := c.db.Collection(c.companyCollection)

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"name":        company.Name,
		"branch":      company.Branch,
		"phone":       company.Phone,
		"mobile":      company.Mobile,
		"description": company.Description,
		"imageURL":    company.ImageURL,
	}}

	_, err := companyCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (c *companyRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Company, error) {
	companyCollection := c.db.Collection(c.companyCollection)

	filter := bson.M{"_id": id}
	var company domain.Company
	if err := companyCollection.FindOne(ctx, filter).Decode(&company); err != nil {
		return nil, err
	}

	return &company, nil
}

func (c *companyRepository) GetByName(ctx context.Context, name string) (*domain.Company, error) {
	companyCollection := c.db.Collection(c.companyCollection)

	filter := bson.M{"name": name}
	var company domain.Company
	if err := companyCollection.FindOne(ctx, filter).Decode(&company); err != nil {
		return nil, err
	}

	return &company, nil
}

func (c *companyRepository) GetAll(ctx context.Context) ([]domain.Company, error) {
	companyCollection := c.db.Collection(c.companyCollection)

	filter := bson.M{}
	cursor, err := companyCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var companies []domain.Company
	for cursor.Next(ctx) {
		var company domain.Company
		if err := cursor.Decode(&company); err != nil {
			return nil, err
		}

		companies = append(companies, company)
	}

	return companies, nil
}
