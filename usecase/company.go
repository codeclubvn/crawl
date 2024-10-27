package usecase

import (
	"context"
	"crawl/domain"
	"crawl/pkg/crawl_data/goquery/trangvangvietnam"
	"crawl/repository"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"time"
)

type ICompanyUseCase interface {
	CreateOne(ctx context.Context, page []string) error
	DeleteOne(ctx context.Context, id primitive.ObjectID) error
	UpdateOne(ctx context.Context, id primitive.ObjectID, company *domain.Company) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Company, error)
	GetByName(ctx context.Context, name string) (*domain.Company, error)
	GetAll(ctx context.Context) ([]domain.Company, error)
}

type companyUseCase struct {
	companyScrap      *trangvangvietnam.ICompanyCrawl
	companyRepository repository.ICompanyRepository
	contextTimeout    time.Duration
}

func NewCompanyUseCase(companyRepository repository.ICompanyRepository, contextTimeout time.Duration) ICompanyUseCase {
	return &companyUseCase{companyRepository: companyRepository, contextTimeout: contextTimeout}
}

func (c *companyUseCase) CreateOne(ctx context.Context, page []string) error {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	currentUrl := page[0]
	companies := *c.companyScrap
	err := companies.GetByTotalPages(currentUrl)
	if err != nil {
		return err
	}

	err = companies.GetByURL(currentUrl)
	if err != nil {
		return err
	}

	err = companies.GetAll(currentUrl)
	if err != nil {
		return err
	}

	company, err := json.Marshal(c.companyScrap) // Chuyển kiểu dữ liệu Ebooks sang JSON
	err = os.WriteFile("output.json", company, 0644)

	return nil
}

func (c *companyUseCase) DeleteOne(ctx context.Context, id primitive.ObjectID) error {
	//TODO implement me
	panic("implement me")
}

func (c *companyUseCase) UpdateOne(ctx context.Context, id primitive.ObjectID, company *domain.Company) error {
	//TODO implement me
	panic("implement me")
}

func (c *companyUseCase) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Company, error) {
	//TODO implement me
	panic("implement me")
}

func (c *companyUseCase) GetByName(ctx context.Context, name string) (*domain.Company, error) {
	//TODO implement me
	panic("implement me")
}

func (c *companyUseCase) GetAll(ctx context.Context) ([]domain.Company, error) {
	//TODO implement me
	panic("implement me")
}