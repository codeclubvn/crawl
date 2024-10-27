package trangvangvietnam

type CompanyCrawl struct{}

type ICompanyCrawl interface {
	GetByURL(url string) error
	GetByTotalPages(url string) error
	GetAll(currentUrl string) error
}

func NewCompaniesCrawl() ICompanyCrawl {
	return &CompanyCrawl{}
}

func (c *CompanyCrawl) GetByURL(url string) error {
	//TODO implement me
	panic("implement me")
}

func (c *CompanyCrawl) GetByTotalPages(url string) error {
	//TODO implement me
	panic("implement me")
}

func (c *CompanyCrawl) GetAll(currentUrl string) error {
	//TODO implement me
	panic("implement me")
}
