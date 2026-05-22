package models

import "time"

type UserStats struct {
	UserID        string `json:"user_id"`
	AcceptedCount int    `json:"accepted_count_current_window"`
	RejectedCount int    `json:"rejected_count_total"`
}

type StatsResponse struct {
	Stats []UserStats `json:"stats"`
}

type ProductSummaryResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	SKU          string    `json:"sku"`
	ImageCount   int       `json:"image_count"`
	VideoCount   int       `json:"video_count"`
	ThumbnailURL string    `json:"thumbnail_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type ProductDetailResponse struct {
	ProductSummaryResponse
	ImageURLs []string `json:"image_urls"`
	VideoURLs []string `json:"video_urls"`
}

// Internal Representation of Product
type InternalProduct struct {
	ProductSummaryResponse
	ImageURLs []string
	VideoURLs []string
}
