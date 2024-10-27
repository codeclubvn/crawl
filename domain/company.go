package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Company struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	Name        string             `bson:"name" json:"name"`
	Branch      string             `bson:"branch" json:"branch"`
	Phone       string             `bson:"phone" json:"phone"`
	Mobile      string             `bson:"mobile" json:"mobile"`
	Description string             `bson:"description" json:"description"`
	ImageURL    string             `bson:"imageURL" json:"imageURL"`
}

type Companies struct {
	TotalPages     int       `json:"totalPages"`
	TotalCompanies int       `json:"totalCompanies"`
	List           []Company `json:"companies"`
}
