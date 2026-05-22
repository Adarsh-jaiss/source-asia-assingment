package controllers

import (
	"net/http"

	"github.com/adarsh-jaiss/assingment/api/models"
	"github.com/adarsh-jaiss/assingment/api/repository"
	"github.com/adarsh-jaiss/assingment/utils"
	"github.com/gin-gonic/gin"
)

type RateLimitController struct {
	repo repository.RateLimitRepo
}

func NewRateLimitController(repo repository.RateLimitRepo) *RateLimitController {
	return &RateLimitController{repo: repo}
}

// HandleRequest godoc
// @Summary Record a request and enforce rate limiting
// @Description Accepts a request and checks if it's within the rate limit. Returns 201 if accepted, 429 if rejected.
// @Tags RateLimit
// @Accept json
// @Produce json
// @Param request body models.RateLimitRequest true "Request Payload"
// @Success 201 {object} map[string]string
// @Failure 400 {object} utils.ErrorResponse
// @Failure 429 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/request [post]
func (c *RateLimitController) HandleRequest(ctx *gin.Context) {
	var req models.RateLimitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.JSONError(ctx, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request payload")
		return
	}

	if req.UserID == "" {
		utils.JSONError(ctx, http.StatusBadRequest, "INVALID_USER_ID", "user_id cannot be empty")
		return
	}

	allowed, err := c.repo.CheckAndRecord(req.UserID)
	if err != nil {
		utils.JSONError(ctx, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	if !allowed {
		utils.JSONError(ctx, http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED", "Rate limit exceeded. Maximum 5 requests per minute.")
		return
	}

	utils.JSONSuccess(ctx, http.StatusCreated, map[string]string{"message": "Request accepted"})
}

// GetStats godoc
// @Summary Retrieve rate limiting statistics
// @Description Returns the accepted requests in the current window and total rejected requests per user
// @Tags RateLimit
// @Produce json
// @Success 200 {object} models.StatsResponse
// @Router /api/v1/stats [get]
func (c *RateLimitController) GetStats(ctx *gin.Context) {
	stats := c.repo.GetStats()
	utils.JSONSuccess(ctx, http.StatusOK, models.StatsResponse{Stats: stats})
}
