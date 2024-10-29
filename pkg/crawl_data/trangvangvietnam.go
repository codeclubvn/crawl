package crawl_data

import (
	"context"
	"crawl/domain"
	"crawl/repository"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

type CompanyScraping struct {
	companies         domain.Companies
	companyRepository repository.ICompanyRepository
}

type ICompanyCrawl interface {
	GetByURL(url string) error
	GetByTotalPages(url string) error
	GetAll(currentUrl string) error
}

func NewCompaniesCrawl(companies domain.Companies, companyRepository repository.ICompanyRepository) ICompanyCrawl {
	return &CompanyScraping{companies: companies, companyRepository: companyRepository}
}

func (c *CompanyScraping) GetByURL(url string) error {
	// Gửi yêu cầu HTTP tới URL
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Tạo tài liệu từ phản hồi
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	// Mutex để bảo vệ truy cập vào c.companies.List trong môi trường đa luồng
	var mu sync.Mutex

	// Duyệt qua từng công ty trong danh sách
	doc.Find(".w-100.h-auto.shadow.rounded-3.bg-white.p-2.mb-3").Each(func(i int, s *goquery.Selection) {
		// Tên công ty
		name := s.Find(".listings_center h2 a").Text()

		// Đường dẫn trang công ty
		companyURL, exists := s.Find(".listings_center h2 a").Attr("href")
		if !exists {
			companyURL = "" // Gán giá trị rỗng nếu không tìm thấy
		}

		// Ngành nghề
		branch := s.Find("span.nganh_listing_txt").Text()

		// Địa chỉ
		address := s.Find(".logo_congty_diachi .fa-location-dot").Parent().Text()

		// Số điện thoại
		phone := s.Find(".listing_dienthoai a").Text()

		// Hotline
		hotline := s.Find("span.fw500 a").Text()

		// Mô tả
		description := s.Find("div.div_textqc > small").Text()

		// URL hình ảnh logo công ty
		imageURL, exists := s.Find(".logo_congty img").Attr("src")
		if !exists {
			imageURL = "" // Gán giá trị rỗng nếu không tìm thấy
		}

		// Tạo đối tượng Company và thêm vào danh sách
		company := domain.Company{
			Name:        name,
			Address:     address,
			Branch:      branch,
			Phone:       phone,
			Mobile:      hotline,
			Description: description,
			ImageURL:    imageURL,
			CompanyURL:  companyURL,
		}

		mu.Lock()
		c.companies.TotalCompanies++
		c.companies.List = append(c.companies.List, company)
		mu.Unlock()
	})

	return nil
}

func (c *CompanyScraping) GetByTotalPages(url string) error {
	// Gửi yêu cầu HTTP đến URL
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Tạo tài liệu từ phản hồi
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	// Tìm liên kết của trang cuối trong phần phân trang
	lastPageLink, exists := doc.Find("div#paging a").Eq(-2).Attr("href") // Lấy phần tử áp chót để tránh nút "Tiếp"
	if !exists || lastPageLink == "#" {                                  // Trường hợp chỉ có 1 trang hoặc không có phân trang
		c.companies.TotalPages = 1
		return nil
	}

	// Tách số trang từ URL cuối
	split := strings.Split(lastPageLink, "?page=")
	if len(split) < 2 {
		return fmt.Errorf("could not parse the total pages from URL: %s", lastPageLink)
	}

	totalPages, err := strconv.Atoi(split[1])
	if err != nil {
		return fmt.Errorf("invalid page number in URL %s: %w", lastPageLink, err)
	}

	// Gán tổng số trang vào `c.companies.TotalPages`
	c.companies.TotalPages = totalPages
	return nil
}

func (c *CompanyScraping) GetAll(currentUrl string) error {
	eg := errgroup.Group{}
	if c.companies.TotalPages > 0 {
		for i := 1; i <= c.companies.TotalPages; i++ {
			// Lưu giá trị của uri cho mỗi lần lặp để tránh bị ghi đè
			uri := fmt.Sprintf("%v?page=%v", currentUrl, i)
			eg.Go(func(uri string) func() error {
				return func() error {
					return c.GetByURL(uri)
				}
			}(uri))
		}

		// Chờ tất cả các goroutines hoàn thành
		if err := eg.Wait(); err != nil {
			return err
		}
	}

	// Tạo dữ liệu để ghi vào file
	companyData := domain.Companies{
		TotalPages:     c.companies.TotalPages,
		TotalCompanies: c.companies.TotalCompanies,
		List:           c.companies.List,
	}

	var egInsert errgroup.Group
	for _, company := range c.companies.List {
		company := company // Tạo bản sao của company để tránh vấn đề với goroutine
		egInsert.Go(func() error {
			err := c.companyRepository.CreateOne(context.Background(), &company)
			return err
		})
	}

	// Chờ tất cả các goroutines hoàn thành
	if err := egInsert.Wait(); err != nil {
		return fmt.Errorf("failed to insert companies into MongoDB: %w", err)
	}

	// Marshal dữ liệu sau khi tất cả goroutines hoàn thành
	companiesData, err := json.Marshal(companyData)
	if err != nil {
		return fmt.Errorf("failed to marshal company data: %w", err)
	}

	// Ghi dữ liệu JSON vào file
	err = os.WriteFile("output.json", companiesData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write company data to file: %w", err)
	}
	return nil
}
