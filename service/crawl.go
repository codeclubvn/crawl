package service

import (
	"crawl/model"
	"crawl/repo"
	"fmt"
	"github.com/gocolly/colly"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Crawl struct {
	repo repo.IRepo
}

func NewCrawl(repo repo.IRepo) *Crawl {
	return &Crawl{
		repo: repo,
	}
}

type IUser interface {
	CrawlYellowPage(baseURL string) (err error)
	GetOne(filter map[string]interface{}) (model.CompanyInfo, error)
	GetList(filter map[string]interface{}, limit int64, skip int64) ([]model.CompanyInfo, error)
}

func (u *Crawl) CrawlYellowPage(baseURL string) (err error) {
	var (
		companies []model.CompanyInfo
		images    []model.Image
		mu        sync.Mutex
	)
	// Tạo channel để xử lý companies
	companiesChan := make(chan []model.CompanyInfo, 100)
	var wg sync.WaitGroup

	// Khởi tạo collector với Async true để chạy parallel
	co := colly.NewCollector(
		colly.Async(true), // Enable async processing
		colly.AllowURLRevisit(),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36"),
	)

	// Tăng số lượng parallel requests
	err = co.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 5, // Tăng số lượng parallel requests
		RandomDelay: 2 * time.Second,
	})
	if err != nil {
		return
	}

	// Worker để xử lý insert data
	wg.Add(1)
	go func() {
		defer wg.Done()
		for companies := range companiesChan {
			err := u.repo.CreateBulk(companies)
			if err != nil {
				fmt.Printf("Error inserting companies: %v\n", err)
			}
		}
	}()

	co.OnHTML("div.w-100.h-auto.shadow.rounded-3.bg-white.p-2.mb-3", func(e *colly.HTMLElement) {
		company := model.CompanyInfo{}

		// Get company name
		company.Name = strings.TrimSpace(e.ChildText("h2.p-1.fs-5.h2.m-0.pt-0.ps-0.text-capitalize a"))

		// Get address
		addressDiv := e.ChildText("div.pt-0.pb-2.ps-3.pe-4 small")
		if addressDiv != "" {
			company.Address = strings.TrimSpace(strings.Split(addressDiv, "Việt Nam")[0] + "Việt Nam")
		}

		// Get phone numbers
		phones := []string{}
		e.ForEach("div.listing_dienthoai a", func(_ int, el *colly.HTMLElement) {
			phone := strings.TrimSpace(el.Text)
			if phone != "" {
				phones = append(phones, phone)
			}
		})
		company.Phone = phones

		// Get hotline
		company.Hotline = strings.TrimSpace(e.ChildText("div.pt-0.pb-2.ps-3.pe-4 span.fw500 a"))

		// Get industry
		company.Industry = strings.TrimSpace(e.ChildText("span.nganh_listing_txt.fw500"))

		// Get description
		company.Description = strings.TrimSpace(e.ChildText("div.div_textqc small"))

		// Get email
		emailEl := e.ChildAttr("div.email_web_section a[href^='mailto:']", "href")
		company.Email = strings.TrimPrefix(emailEl, "mailto:")

		// Get website
		websiteEl := e.ChildAttr("div.email_web_section a[rel='nofollow']", "href")
		company.Website = websiteEl

		// Get images
		e.ForEach("div.div_showimages img", func(_ int, el *colly.HTMLElement) {
			img := model.Image{
				URL:         el.Attr("src"),
				Title:       el.Attr("title"),
				Description: el.Attr("alt"),
			}
			images = append(images, img)
		})
		company.Images = images

		mu.Lock()
		companies = append(companies, company)
		mu.Unlock()
	})

	// Thêm callback này để insert data sau mỗi trang
	co.OnScraped(func(r *colly.Response) {
		mu.Lock()
		if len(companies) > 0 {
			companiesChan <- companies
			companies = []model.CompanyInfo{}
		}
		mu.Unlock()
	})

	// Xử lý pagination và chuyển trang
	co.OnHTML("div#paging", func(e *colly.HTMLElement) {
		// Lấy trang cuối cùng
		lastPage := 1
		e.ForEach("a", func(_ int, el *colly.HTMLElement) {
			if pageNum, err := strconv.Atoi(el.Text); err == nil && pageNum > lastPage {
				lastPage = pageNum
			}
		})

		// Lấy trang hiện tại từ URL
		currentPage := 1
		if pageParam := e.Request.URL.Query().Get("page"); pageParam != "" {
			if page, err := strconv.Atoi(pageParam); err == nil {
				currentPage = page
			}
		}

		// Nếu chưa phải trang cuối, crawl trang tiếp theo
		if currentPage < lastPage {
			nextPage := currentPage + 1
			nextURL := fmt.Sprintf("%s?page=%d", baseURL, nextPage)
			time.Sleep(2 * time.Second) // Thêm delay để tránh request quá nhanh
			co.Visit(nextURL)
		}
	})

	// Logging cho debugging
	co.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Thêm error handling
	co.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Thực hiện visit
	err = co.Visit(baseURL)
	if err != nil {
		fmt.Println("Error visiting URL:", err)
		return
	}

	// Đợi tất cả requests hoàn thành
	co.Wait()

	// Đóng channel và đợi worker hoàn thành
	close(companiesChan)
	wg.Wait()

	return
}

func (u *Crawl) GetOne(filter map[string]interface{}) (model.CompanyInfo, error) {
	return u.repo.GetOne(filter)
}

func (u *Crawl) GetList(filter map[string]interface{}, limit int64, skip int64) ([]model.CompanyInfo, error) {
	return u.repo.GetList(filter, limit, skip)
}
