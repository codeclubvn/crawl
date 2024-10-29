package domain

type Company struct {
	Name        string `bson:"name" json:"name"`
	Address     string `bson:"address" json:"address"`
	Branch      string `bson:"branch" json:"branch"`
	Phone       string `bson:"phone" json:"phone"`
	Mobile      string `bson:"mobile" json:"mobile"`
	Description string `bson:"description" json:"description"`
	ImageURL    string `bson:"imageURL" json:"imageURL"`
	CompanyURL  string `bson:"companyURL" json:"companyURL"`
}

type Companies struct {
	TotalPages     int       `json:"totalPages"`
	TotalCompanies int       `json:"totalCompanies"`
	List           []Company `json:"companies"`
}
