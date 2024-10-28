package model

type CompanyInfo struct {
	Name        string   // Tên công ty
	Address     string   // Địa chỉ
	Phone       []string // Số điện thoại
	Hotline     string   // Số hotline
	Email       string   // Email công ty
	Website     string   // Website công ty
	Description string   // Mô tả công ty
	Industry    string   // Ngành nghề
	Images      []Image  // Hình ảnh về công ty và dịch vụ
}

type Image struct {
	URL         string
	Title       string
	Description string
}

type GetOneInput struct {
	Filter map[string]interface{} `json:"filter"`
}

type GetListInput struct {
	Filter map[string]interface{} `json:"filter"`
	Limit  int64                  `json:"limit"`
	Skip   int64                  `json:"skip"`
}
