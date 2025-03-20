package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

type Server struct {
	service service
}

func New(service service) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) NewRoutes() http.Handler {
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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.GET("/api/v1/products/", s.getProducts)
	router.GET("/api/v1/products/:voucher_code", s.getProductsWithVoucher)
	router.GET("/api/v1/product/:product_id", s.getProduct)
	router.POST("/api/v1/product/subscribe/", s.subscribe)
	router.GET("/api/v1/subscription/:subscription_id", s.getSubscription)
	router.POST("/api/v1/subscription/:subscription_id/manage", s.manageSubscription)

	return router
}
