package rest

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) getProducts(c *gin.Context) {
	ctx := context.Background()

	products, err := s.service.FindProducts(ctx)
	if err != nil {
		log.Printf("Error finding products: %v", err)
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Products not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (s *Server) getProductsWithVoucher(c *gin.Context) {
	ctx := context.Background()

	voucherCode := c.Param("voucher_code")
	products, err := s.service.FindProductsWithVoucher(ctx, voucherCode)
	if err != nil {
		log.Printf("Error finding products with voucher: %v", err)
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Products with voucher not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (s *Server) getProduct(c *gin.Context) {
	ctx := context.Background()

	productID := c.Param("product_id")
	product, err := s.service.FindProduct(ctx, productID)
	if err != nil {
		log.Printf("Error finding product with ID %s: %v", productID, err)
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Product not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (s *Server) subscribe(c *gin.Context) {
	ctx := context.Background()

	var request struct {
		UserID      string `json:"user_id" binding:"required"`
		ProductID   string `json:"product_id" binding:"required"`
		VoucherCode string `json:"voucher_code"`
		TrialPeriod bool   `json:"trial_period"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Validation error: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error", "details": err.Error()})
		return
	}

	if request.UserID == "" || request.ProductID == "" {
		log.Println("Validation error: missing user_id or product_id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error: missing user_id or product_id"})
		return
	}

	subscriptionID, err := s.service.Subscribe(
		ctx,
		request.UserID,
		request.ProductID,
		request.VoucherCode,
		request.TrialPeriod,
	)
	if err != nil {
		log.Printf("Error subscribing user %s to product %s: %v", request.UserID, request.ProductID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal error",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"subscription_id": subscriptionID})
}

func (s *Server) getSubscription(c *gin.Context) {
	ctx := context.Background()

	subscriptionID := c.Param("subscription_id")
	subscription, err := s.service.FindSubscription(ctx, subscriptionID)
	if err != nil {
		log.Printf("Error finding subscription with ID %s: %v", subscriptionID, err)
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Subscription not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

func (s *Server) manageSubscription(c *gin.Context) {
	ctx := context.Background()

	subscriptionID := c.Param("subscription_id")

	var request struct {
		Action string `json:"action" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action"})
		return
	}

	switch request.Action {
	case "pause":
		err := s.service.PauseSubscription(ctx, subscriptionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Subscription paused"})

	case "unpause":
		err := s.service.UnpauseSubscription(ctx, subscriptionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Subscription unpaused"})

	case "cancel":
		err := s.service.CancelSubscription(ctx, subscriptionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Subscription canceled"})

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action"})
	}
}
