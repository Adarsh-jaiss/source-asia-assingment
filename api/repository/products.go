package repository

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"github.com/adarsh-jaiss/assingment/api/models"
)

var (
	ErrDuplicateSKU = errors.New("duplicate sku")
	ErrProductNotFound = errors.New("product not found")
)

type ProductRepo interface {
	CreateProduct(req *models.ProductCreateRequest) (*models.InternalProduct, error)
	GetProducts(limit, offset int) ([]models.ProductSummaryResponse, error)
	GetProductByID(id string) (*models.InternalProduct, error)
	AppendMedia(id string, req *models.MediaAppendRequest) (*models.InternalProduct, error)
}

type productRepoImpl struct {
	mu           sync.RWMutex
	productsByID map[string]*models.InternalProduct
	productIDs   []string // For stable ordering and pagination
	skus         map[string]bool
}

func NewProductRepo() ProductRepo {
	return &productRepoImpl{
		productsByID: make(map[string]*models.InternalProduct),
		productIDs:   make([]string, 0),
		skus:         make(map[string]bool),
	}
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (r *productRepoImpl) CreateProduct(req *models.ProductCreateRequest) (*models.InternalProduct, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.skus[req.SKU] {
		return nil, ErrDuplicateSKU
	}

	id := generateID()
	now := time.Now()

	thumbnailURL := ""
	if len(req.ImageURLs) > 0 {
		thumbnailURL = req.ImageURLs[0]
	}

	product := &models.InternalProduct{
		ProductSummaryResponse: models.ProductSummaryResponse{
			ID:           id,
			Name:         req.Name,
			SKU:          req.SKU,
			ImageCount:   len(req.ImageURLs),
			VideoCount:   len(req.VideoURLs),
			ThumbnailURL: thumbnailURL,
			CreatedAt:    now,
		},
		ImageURLs: req.ImageURLs,
		VideoURLs: req.VideoURLs,
	}

	r.productsByID[id] = product
	r.productIDs = append(r.productIDs, id)
	r.skus[req.SKU] = true

	return product, nil
}

func (r *productRepoImpl) GetProducts(limit, offset int) ([]models.ProductSummaryResponse, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	total := len(r.productIDs)
	if offset > total {
		return []models.ProductSummaryResponse{}, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}

	result := make([]models.ProductSummaryResponse, 0, end-offset)
	for i := offset; i < end; i++ {
		id := r.productIDs[i]
		product := r.productsByID[id]
		result = append(result, product.ProductSummaryResponse)
	}

	return result, nil
}

func (r *productRepoImpl) GetProductByID(id string) (*models.InternalProduct, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	product, exists := r.productsByID[id]
	if !exists {
		return nil, ErrProductNotFound
	}
	return product, nil
}

func (r *productRepoImpl) AppendMedia(id string, req *models.MediaAppendRequest) (*models.InternalProduct, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	product, exists := r.productsByID[id]
	if !exists {
		return nil, ErrProductNotFound
	}

	if len(req.ImageURLs) > 0 {
		product.ImageURLs = append(product.ImageURLs, req.ImageURLs...)
		product.ImageCount = len(product.ImageURLs)
		if product.ThumbnailURL == "" {
			product.ThumbnailURL = product.ImageURLs[0]
		}
	}

	if len(req.VideoURLs) > 0 {
		product.VideoURLs = append(product.VideoURLs, req.VideoURLs...)
		product.VideoCount = len(product.VideoURLs)
	}

	return product, nil
}