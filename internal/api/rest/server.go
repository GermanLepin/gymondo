package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	"net/http"
)

type Service struct {
	productService productService
}

func New(productService productService) *Service {
	return &Service{
		productService: productService,
	}
}

func (s *Service) NewRoutes() http.Handler {
	router := gin.Default()

	// specify who is allowed to connect
	router.Use(func(c *gin.Context) {
		corsMiddleware := cors.New(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		})
		corsMiddleware.HandlerFunc(c.Writer, c.Request)
		c.Next()
	})

	router.GET("/api/products/:product_id", s.getProduct)
	router.GET("/api/products", s.getProducts)
	//router.POST("/subscriptions", buyProduct)

	return router
}
