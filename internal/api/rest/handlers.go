package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// @Summary Get all products
// @Description Retrieves a list of all available products. This endpoint provides information about the products that users can subscribe to.
// @Tags Products
// @Produce json
// @Success 200 {array} model.Product
// @Failure 404 {object} ErrorResponse "Products not found"
// @Router /api/v1/products [get]
func (s *Server) getProducts(c *gin.Context) {
	ctx := context.Background()

	products, err := s.service.FindProducts(ctx)
	if err != nil {
		log.Printf("Error finding products: %v", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Products not found",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, products)
}

// @Summary Get all products with a voucher
// @Description Fetches details of a specific product associated with a given voucher code. The voucher code is used to apply discounts or offers to the product.
// @Tags Products
// @Produce json
// @Param voucher_code path string true "Voucher Code"
// @Success 200 {array} model.Product
// @Failure 404 {object} ErrorResponse "Products with voucher not found"
// @Router /api/v1/products/{voucher_code} [get]
func (s *Server) getProductsWithVoucher(c *gin.Context) {
	ctx := context.Background()

	voucherCode := c.Param("voucher_code")
	products, err := s.service.FindProductsWithVoucher(ctx, voucherCode)
	if err != nil {
		log.Printf("Error finding products with voucher: %v", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Products with voucher not found",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, products)
}

// @Summary Get a specific product
// @Description  Retrieves detailed information about a specific product using the unique product_id. This includes pricing, description, and other attributes.
// @Tags Product
// @Produce json
// @Param product_id path string true "Product ID"
// @Success 200 {object} model.Product
// @Failure 404 {object} ErrorResponse "Product not found"
// @Router /api/v1/product/{product_id} [get]
func (s *Server) getProduct(c *gin.Context) {
	ctx := context.Background()

	productID := c.Param("product_id")
	product, err := s.service.FindProduct(ctx, productID)
	if err != nil {
		log.Printf("Error finding product with ID %s: %v", productID, err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Product not found",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

type SubscriptionRequest struct {
	UserID      string `json:"user_id" binding:"required"`
	ProductID   string `json:"product_id" binding:"required"`
	VoucherCode string `json:"voucher_code,omitempty"`
	TrialPeriod bool   `json:"trial_period"`
}

type SubscriptionResponse struct {
	SubscriptionID string `json:"subscription_id"`
	Message        string `json:"message"`
}

// @Summary Subscribe to a product
// @Description Allows users to subscribe to a product. This endpoint creates a new subscription for a user, including selecting a product and setting the subscription parameters (e.g., trial period, voucher code).
// @Tags Product
// @Accept json
// @Produce json
// @Param request body SubscriptionRequest true "Subscription Request"
// @Success 200 {object} SubscriptionResponse
// @Failure 400 {object} ErrorResponse "Validation error"
// @Failure 500 {object} ErrorResponse "Internal error"
// @Router /api/v1/product/subscribe [post]
func (s *Server) subscribe(c *gin.Context) {
	ctx := context.Background()

	var request SubscriptionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Validation error: ", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation error",
			Details: err.Error(),
		})
		return
	}

	subscriptionID, err := s.service.Subscribe(ctx, request.UserID, request.ProductID, request.VoucherCode, request.TrialPeriod)
	if err != nil {
		log.Printf("Error subscribing user %s to product %s: %v", request.UserID, request.ProductID, err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal error",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SubscriptionResponse{
		SubscriptionID: subscriptionID,
		Message:        "Subscription created successfully",
	})
}

// @Summary Get subscription details
// @Description Provides details of an active subscription. The subscription_id is used to fetch information about a specific subscription, such as its status, start date, end date, and other relevant information.
// @Tags Subscription
// @Produce json
// @Param subscription_id path string true "Subscription ID"
// @Success 200 {object} model.Subscription
// @Failure 404 {object} ErrorResponse "Subscription not found"
// @Router /api/v1/subscription/{subscription_id} [get]
func (s *Server) getSubscription(c *gin.Context) {
	ctx := context.Background()

	subscriptionID := c.Param("subscription_id")
	subscription, err := s.service.FindSubscription(ctx, subscriptionID)
	if err != nil {
		log.Printf("Error finding subscription with ID %s: %v", subscriptionID, err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Subscription not found",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

type ManageSubscriptionRequest struct {
	Action string `json:"action" binding:"required"`
}

type ManageSubscriptionResponse struct {
	SubscriptionID string `json:"subscription_id"`
	Message        string `json:"message"`
}

// @Summary Manage subscription
// @Description Manages an existing subscription. This endpoint allows users to update or modify their subscription, such as pausing, canceling, or changing other settings related to the subscription.
// @Tags Subscription
// @Accept json
// @Produce json
// @Param subscription_id path string true "Subscription ID"
// @Param request body ManageSubscriptionRequest true "Manage Action"
// @Success 200 {object} ManageSubscriptionResponse
// @Failure 400 {object} ErrorResponse "Invalid action"
// @Failure 500 {object} ErrorResponse "Internal error"
// @Router /api/v1/subscription/{subscription_id}/manage [post]
func (s *Server) manageSubscription(c *gin.Context) {
	ctx := context.Background()
	subscriptionID := c.Param("subscription_id")

	var request ManageSubscriptionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid action",
			Details: err.Error(),
		})
		return
	}

	switch request.Action {
	case "pause":
		err := s.service.PauseSubscription(ctx, subscriptionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Failed to pause subscription",
				Details: fmt.Sprintf("Error pausing subscription: %v", err),
			})
			return
		}

		c.JSON(http.StatusOK, SubscriptionResponse{
			SubscriptionID: subscriptionID,
			Message:        "Subscription paused",
		})
	case "unpause":
		err := s.service.UnpauseSubscription(ctx, subscriptionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Failed to unpause subscription",
				Details: fmt.Sprintf("Error unpausing subscription: %v", err),
			})
			return
		}

		c.JSON(http.StatusOK, SubscriptionResponse{
			SubscriptionID: subscriptionID,
			Message:        "Subscription unpaused",
		})
	case "cancel":
		err := s.service.CancelSubscription(ctx, subscriptionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Failed to cancel subscription",
				Details: fmt.Sprintf("Error canceling subscription: %v", err),
			})
			return
		}

		c.JSON(http.StatusOK, SubscriptionResponse{
			SubscriptionID: subscriptionID,
			Message:        "Subscription canceled",
		})
	default:
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid action",
			Details: fmt.Sprintf("Action '%s' is not supported", request.Action),
		})
	}
}
