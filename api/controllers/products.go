package controllers

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/adarsh-jaiss/assingment/api/models"
	"github.com/adarsh-jaiss/assingment/api/repository"
	"github.com/adarsh-jaiss/assingment/utils"
	"github.com/gin-gonic/gin"
)

type ProductController struct {
	repo repository.ProductRepo
}

func NewProductController(repo repository.ProductRepo) *ProductController {
	return &ProductController{repo: repo}
}

func validateURLs(urls []string) bool {
	for _, u := range urls {
		parsedURL, err := url.ParseRequestURI(u)
		if err != nil {
			return false
		}
		if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
			return false
		}
		if len(u) > 2048 {
			return false
		}
	}
	return true
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Creates a product with name, sku, and media URLs. Validates that max 20 URLs are passed per type.
// @Tags Products
// @Accept json
// @Produce json
// @Param request body models.ProductCreateRequest true "Product Request Payload"
// @Success 201 {object} models.ProductDetailResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 409 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/products [post]
func (c *ProductController) CreateProduct(ctx *gin.Context) {
	var req models.ProductCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.JSONError(ctx, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.SKU = strings.TrimSpace(req.SKU)

	if req.Name == "" || req.SKU == "" {
		utils.JSONError(ctx, http.StatusBadRequest, "INVALID_INPUT", "name and sku cannot be empty")
		return
	}

	if len(req.ImageURLs) > 20 || len(req.VideoURLs) > 20 {
		utils.JSONError(ctx, http.StatusBadRequest, "LIMIT_EXCEEDED", "Maximum 20 URLs allowed per array")
		return
	}

	if !validateURLs(req.ImageURLs) || !validateURLs(req.VideoURLs) {
		utils.JSONError(ctx, http.StatusBadRequest, "INVALID_URL", "Invalid URL format or length")
		return
	}

	product, err := c.repo.CreateProduct(&req)
	if err != nil {
		if err == repository.ErrDuplicateSKU {
			utils.JSONError(ctx, http.StatusConflict, "DUPLICATE_SKU", "SKU already exists")
			return
		}
		utils.JSONError(ctx, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	// For detail view on creation, we can return ProductDetailResponse
	detailResponse := models.ProductDetailResponse{
		ProductSummaryResponse: product.ProductSummaryResponse,
		ImageURLs:              product.ImageURLs,
		VideoURLs:              product.VideoURLs,
	}

	utils.JSONSuccess(ctx, http.StatusCreated, detailResponse)
}

// GetProducts godoc
// @Summary List products
// @Description Retrieves a paginated list of products. Does not include full image or video URL arrays.
// @Tags Products
// @Produce json
// @Param limit query int false "Pagination limit" default(20)
// @Param offset query int false "Pagination offset" default(0)
// @Success 200 {array} models.ProductSummaryResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/products [get]
func (c *ProductController) GetProducts(ctx *gin.Context) {
	limitStr := ctx.DefaultQuery("limit", "20")
	offsetStr := ctx.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100 // max limit
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	products, err := c.repo.GetProducts(limit, offset)
	if err != nil {
		utils.JSONError(ctx, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	utils.JSONSuccess(ctx, http.StatusOK, products)
}

// GetProductByID godoc
// @Summary Get a product by ID
// @Description Retrieves full product details including image and video URL arrays.
// @Tags Products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} models.ProductDetailResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/products/{id} [get]
func (c *ProductController) GetProductByID(ctx *gin.Context) {
	id := ctx.Param("id")

	product, err := c.repo.GetProductByID(id)
	if err != nil {
		if err == repository.ErrProductNotFound {
			utils.JSONError(ctx, http.StatusNotFound, "NOT_FOUND", "Product not found")
			return
		}
		utils.JSONError(ctx, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	detailResponse := models.ProductDetailResponse{
		ProductSummaryResponse: product.ProductSummaryResponse,
		ImageURLs:              product.ImageURLs,
		VideoURLs:              product.VideoURLs,
	}

	utils.JSONSuccess(ctx, http.StatusOK, detailResponse)
}

// AppendMedia godoc
// @Summary Append media to a product
// @Description Appends image and video URLs to an existing product. Validates max 20 URLs per array.
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body models.MediaAppendRequest true "Media Append Payload"
// @Success 200 {object} models.ProductDetailResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/products/{id}/media [post]
func (c *ProductController) AppendMedia(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.MediaAppendRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.JSONError(ctx, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if len(req.ImageURLs) == 0 && len(req.VideoURLs) == 0 {
		utils.JSONError(ctx, http.StatusBadRequest, "INVALID_INPUT", "At least one of image_urls or video_urls is required")
		return
	}

	if len(req.ImageURLs) > 20 || len(req.VideoURLs) > 20 {
		utils.JSONError(ctx, http.StatusBadRequest, "LIMIT_EXCEEDED", "Maximum 20 URLs allowed per array")
		return
	}

	if !validateURLs(req.ImageURLs) || !validateURLs(req.VideoURLs) {
		utils.JSONError(ctx, http.StatusBadRequest, "INVALID_URL", "Invalid URL format or length")
		return
	}

	product, err := c.repo.AppendMedia(id, &req)
	if err != nil {
		if err == repository.ErrProductNotFound {
			utils.JSONError(ctx, http.StatusNotFound, "NOT_FOUND", "Product not found")
			return
		}
		utils.JSONError(ctx, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	detailResponse := models.ProductDetailResponse{
		ProductSummaryResponse: product.ProductSummaryResponse,
		ImageURLs:              product.ImageURLs,
		VideoURLs:              product.VideoURLs,
	}

	utils.JSONSuccess(ctx, http.StatusOK, detailResponse)
}