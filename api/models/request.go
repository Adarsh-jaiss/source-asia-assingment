package models

type RateLimitRequest struct {
	UserID  string      `json:"user_id" binding:"required"`
	Payload interface{} `json:"payload" binding:"required"`
}

type ProductCreateRequest struct {
	Name      string   `json:"name" binding:"required"`
	SKU       string   `json:"sku" binding:"required"`
	ImageURLs []string `json:"image_urls"`
	VideoURLs []string `json:"video_urls"`
}

type MediaAppendRequest struct {
	ImageURLs []string `json:"image_urls"`
	VideoURLs []string `json:"video_urls"`
}
