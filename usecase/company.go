package usecase

import (
	"context"
	"crawl/domain"
	"crawl/pkg/crawl_data/goquery/trangvangvietnam"
	"crawl/repository"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	companyScrap      trangvangvietnam.ICompanyCrawl
	companyRepository repository.ICompanyRepository
	contextTimeout    time.Duration
}

func NewCompanyUseCase(companyRepository repository.ICompanyRepository, contextTimeout time.Duration,
	companyScrap trangvangvietnam.ICompanyCrawl) ICompanyUseCase {
	return &companyUseCase{companyRepository: companyRepository, contextTimeout: contextTimeout, companyScrap: companyScrap}
}

func (c *companyUseCase) CreateOne(ctx context.Context, page []string) error {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	currentUrl := page[0]

	// Lấy tổng số trang nếu có hàm GetByTotalPages
	if err := c.companyScrap.GetByTotalPages(currentUrl); err != nil {
		return errors.Wrap(err, "failed to get total pages")
	}

	// Lấy thông tin từng công ty từ URL
	if err := c.companyScrap.GetByURL(currentUrl); err != nil {
		return errors.Wrap(err, "failed to get company info by URL")
	}

	// Lấy toàn bộ dữ liệu cần thiết
	if err := c.companyScrap.GetAll(currentUrl); err != nil {
		return errors.Wrap(err, "failed to get all company info")
	}

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
