package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Service) getProduct(c *gin.Context) {
	ctx := context.Background()

	productID := c.Param("product_id")
	fmt.Println(productID)
	product, err := s.productService.FindOne(ctx, productID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (s *Service) getProducts(c *gin.Context) {
	ctx := context.Background()

	products, err := s.productService.FindAll(ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "products not found"})
		return
	}

	c.JSON(http.StatusOK, products)
}
