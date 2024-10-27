package trangvangvietnam

import (
	"crawl/domain"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/sync/errgroup"
	"strconv"
	"strings"
)

type CompanyScraping struct {
	companies domain.Companies
}

type ICompanyCrawl interface {
	GetByURL(url string) error
	GetByTotalPages(url string) error
	GetAll(currentUrl string) error
}

func NewCompaniesCrawl(companies domain.Companies) ICompanyCrawl {
	return &CompanyScraping{companies: companies}
}

func (c *CompanyScraping) GetByURL(url string) error {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return err
	}

	doc.Find(".div_list_city").Each(func(i int, s *goquery.Selection) {
		name, _ := s.Find(".listings_center a").Attr("title")
		imageURL, _ := s.Find(".listings_center a").Attr("href")
		phone, _ := s.Find(".listing_dienthoai i a").Attr("title")
		branch, _ := s.Find("span.nganh_listing_txt").Attr("title")
		mobile, _ := s.Find("span.fw500 a").Attr("title")
		description, _ := s.Find("div.div_textqc > small").Attr("")
		company := domain.Company{
			Name:        name,
			Branch:      branch,
			Phone:       phone,
			Mobile:      mobile,
			Description: description,
			ImageURL:    imageURL,
		}

		c.companies.TotalCompanies++
		c.companies.List = append(c.companies.List, company)
	})

	return nil
}

func (c *CompanyScraping) GetByTotalPages(url string) error {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return err
	}
	lastPageLink, _ := doc.Find("paging").Attr("href") // Đọc dữ liệu từ thẻ a của ul.pagination
	if lastPageLink == "javascript:void();" {          // Trường hợp chỉ có 1 page thì sẽ không có url
		c.companies.TotalPages = 1
		return nil
	}
	split := strings.Split(lastPageLink, "?page=")
	totalPages, _ := strconv.Atoi(split[1])
	c.companies.TotalPages = totalPages
	return nil
}

func (c *CompanyScraping) GetAll(currentUrl string) error {
	eg := errgroup.Group{}
	if c.companies.TotalPages > 0 {
		for i := 1; i <= c.companies.TotalPages; i++ {
			uri := fmt.Sprintf("%v?page=%v", currentUrl, i)
			// https://golang.org/doc/faq#closures_and_goroutines
			eg.Go(func() error {
				err := c.GetByURL(uri)
				if err != nil {
					return err
				}
				return nil
			})
		}
		if err := eg.Wait(); err != nil {
			return err
		}
	}
	
	return nil
}
