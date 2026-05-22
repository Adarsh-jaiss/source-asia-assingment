package tests

import (
	"fmt"
	"testing"

	"github.com/adarsh-jaiss/assingment/api/models"
	"github.com/adarsh-jaiss/assingment/api/repository"
)

func TestProductRepo_CRUD(t *testing.T) {
	repo := repository.NewProductRepo()

	// 1. Create a product
	createReq := &models.ProductCreateRequest{
		Name:      "Test Product 1",
		SKU:       "SKU-100",
		ImageURLs: []string{"http://example.com/img1.jpg", "http://example.com/img2.jpg"},
		VideoURLs: []string{"http://example.com/video1.mp4"},
	}

	product, err := repo.CreateProduct(createReq)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	if product.Name != createReq.Name || product.SKU != createReq.SKU {
		t.Errorf("Product details mismatch")
	}

	if product.ImageCount != 2 || product.VideoCount != 1 {
		t.Errorf("Expected image/video counts to be 2 and 1, got %d and %d", product.ImageCount, product.VideoCount)
	}

	if product.ThumbnailURL != "http://example.com/img1.jpg" {
		t.Errorf("Expected thumbnail URL to be first image, got %s", product.ThumbnailURL)
	}

	// 2. Reject duplicate SKU
	_, err = repo.CreateProduct(createReq)
	if err != repository.ErrDuplicateSKU {
		t.Errorf("Expected ErrDuplicateSKU, got %v", err)
	}

	// 3. Get product by ID
	fetched, err := repo.GetProductByID(product.ID)
	if err != nil {
		t.Fatalf("Failed to fetch product by ID: %v", err)
	}
	if fetched.ID != product.ID {
		t.Errorf("Fetched product ID mismatch")
	}

	// 4. Append media
	appendReq := &models.MediaAppendRequest{
		ImageURLs: []string{"http://example.com/img3.jpg"},
		VideoURLs: []string{"http://example.com/video2.mp4"},
	}

	updated, err := repo.AppendMedia(product.ID, appendReq)
	if err != nil {
		t.Fatalf("Failed to append media: %v", err)
	}

	if updated.ImageCount != 3 || updated.VideoCount != 2 {
		t.Errorf("Expected updated image/video counts to be 3 and 2, got %d and %d", updated.ImageCount, updated.VideoCount)
	}

	if len(updated.ImageURLs) != 3 || updated.ImageURLs[2] != "http://example.com/img3.jpg" {
		t.Errorf("Image URL not appended correctly")
	}
}

func TestProductRepo_PaginationAndPerformance(t *testing.T) {
	repo := repository.NewProductRepo()

	// Seed 50 products
	for i := 1; i <= 50; i++ {
		sku := fmt.Sprintf("SKU-%03d", i)
		name := fmt.Sprintf("Product %d", i)
		_, err := repo.CreateProduct(&models.ProductCreateRequest{
			Name:      name,
			SKU:       sku,
			ImageURLs: []string{"http://example.com/img.jpg"},
		})
		if err != nil {
			t.Fatalf("Failed to seed product %d: %v", i, err)
		}
	}

	// Fetch page 1 (limit 20, offset 0)
	page1, err := repo.GetProducts(20, 0)
	if err != nil {
		t.Fatalf("Failed to get page 1: %v", err)
	}
	if len(page1) != 20 {
		t.Errorf("Expected 20 products, got %d", len(page1))
	}
	if page1[0].SKU != "SKU-001" || page1[19].SKU != "SKU-020" {
		t.Errorf("Ordering or content mismatch in page 1")
	}

	// Fetch page 2 (limit 20, offset 20)
	page2, err := repo.GetProducts(20, 20)
	if err != nil {
		t.Fatalf("Failed to get page 2: %v", err)
	}
	if len(page2) != 20 {
		t.Errorf("Expected 20 products, got %d", len(page2))
	}
	if page2[0].SKU != "SKU-021" || page2[19].SKU != "SKU-040" {
		t.Errorf("Ordering or content mismatch in page 2")
	}

	// Fetch page 3 (limit 20, offset 40 - only 10 items left)
	page3, err := repo.GetProducts(20, 40)
	if err != nil {
		t.Fatalf("Failed to get page 3: %v", err)
	}
	if len(page3) != 10 {
		t.Errorf("Expected 10 products, got %d", len(page3))
	}
	if page3[0].SKU != "SKU-041" || page3[9].SKU != "SKU-050" {
		t.Errorf("Ordering or content mismatch in page 3")
	}
}
