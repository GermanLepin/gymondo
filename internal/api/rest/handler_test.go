package rest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gymondo/internal/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func performRequest(r *gin.Engine, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func performPostRequest(r *gin.Engine, path string, body string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("POST", path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func Test_GetProducts(t *testing.T) {
	t.Parallel()

	t.Run("successful test", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := NewMockservice(ctrl)
		server := &Server{service: mockService}

		mockService.EXPECT().FindProducts(gomock.Any()).Return([]model.Product{
			{ID: uuid.New(), Name: "Product 1", Price: 100},
			{ID: uuid.New(), Name: "Product 2", Price: 200},
		}, nil)

		r := gin.Default()
		r.GET("/products", server.getProducts)

		w := performRequest(r, "GET", "/products")
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Product 1")
		assert.Contains(t, w.Body.String(), "Product 2")
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := NewMockservice(ctrl)
		server := &Server{service: mockService}

		mockService.EXPECT().FindProducts(gomock.Any()).Return(
			[]model.Product{}, fmt.Errorf("failed to scan product row"),
		)

		r := gin.Default()
		r.GET("/products", server.getProducts)
		w := performRequest(r, "GET", "/products")
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func Test_GetProductsWithVoucher(t *testing.T) {
	t.Parallel()

	t.Run("successful test", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := NewMockservice(ctrl)
		server := &Server{service: mockService}

		voucherCode := "DISCOUNT10"
		mockService.EXPECT().FindProductsWithVoucher(gomock.Any(), voucherCode).Return([]model.Product{
			{ID: uuid.New(), Name: "Product 1", Price: 90},
		}, nil)

		r := gin.Default()
		r.GET("/products/:voucher_code", server.getProductsWithVoucher)
		w := performRequest(r, "GET", "/products/DISCOUNT10")
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Product 1")
	})

	t.Run("products not found", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := NewMockservice(ctrl)
		server := &Server{service: mockService}

		mockService.EXPECT().FindProductsWithVoucher(gomock.Any(), "INVALID").Return(
			[]model.Product{}, fmt.Errorf("products not found"),
		)

		r := gin.Default()
		r.GET("/products/:voucher_code", server.getProductsWithVoucher)
		w := performRequest(r, "GET", "/products/INVALID")
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("404", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := NewMockservice(ctrl)
		server := &Server{service: mockService}

		r := gin.Default()
		r.GET("/products/:voucher_code", server.getProductsWithVoucher)
		w := performRequest(r, "GET", "/products/")
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func Test_GetProduct(t *testing.T) {
	t.Parallel()

	t.Run("successful test", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := NewMockservice(ctrl)
		server := &Server{service: mockService}

		productID := uuid.New()
		expectedProduct := model.Product{ID: productID, Name: "Test Product", Price: 150}

		mockService.EXPECT().FindProduct(gomock.Any(), productID.String()).Return(expectedProduct, nil)

		r := gin.Default()
		r.GET("/api/product/:product_id", server.getProduct)

		w := performRequest(r, "GET", "/api/product/"+productID.String())
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), expectedProduct.Name)
	})

	t.Run("product not found", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := NewMockservice(ctrl)
		server := &Server{service: mockService}

		productID := uuid.New().String()
		mockService.EXPECT().FindProduct(gomock.Any(), productID).Return(model.Product{}, fmt.Errorf("product not found"))

		r := gin.Default()
		r.GET("/api/product/:product_id", server.getProduct)

		w := performRequest(r, "GET", "/api/product/"+productID)
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Product not found")
	})

	t.Run("404", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := NewMockservice(ctrl)
		server := &Server{service: mockService}

		r := gin.Default()
		r.GET("/api/product/:product_id", server.getProduct)

		w := performRequest(r, "GET", "/api/product/")
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func Test_Subscribe(t *testing.T) {
	t.Parallel()

	t.Run("successful subscription", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := NewMockservice(ctrl)
		server := &Server{service: mockService}

		requestBody := `{
							"user_id": "123", 
							"product_id": "456", 
							"voucher_code": "ABC123", 
							"trial_period": true
						}`
		expectedSubscriptionID := uuid.New().String()

		mockService.EXPECT().Subscribe(gomock.Any(), "123", "456", "ABC123", true).Return(expectedSubscriptionID, nil)

		r := gin.Default()
		r.POST("/api/subscribe", server.subscribe)

		w := performPostRequest(r, "/api/subscribe", requestBody)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), expectedSubscriptionID)
	})

	t.Run("missing required fields", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := NewMockservice(ctrl)
		server := &Server{service: mockService}

		invalidRequestBody := `{"user_id": "", "product_id": ""}`

		r := gin.Default()
		r.POST("/api/subscribe", server.subscribe)

		w := performPostRequest(r, "/api/subscribe", invalidRequestBody)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Validation error")
	})

	t.Run("internal service error", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := NewMockservice(ctrl)
		server := &Server{service: mockService}

		requestBody := `{
							"user_id": "123", 
							"product_id": "456", 
							"voucher_code": "ABC123", 
							"trial_period": true
						}`

		mockService.EXPECT().Subscribe(gomock.Any(), "123", "456", "ABC123", true).
			Return("", fmt.Errorf("internal service error"))

		r := gin.Default()
		r.POST("/api/subscribe", server.subscribe)

		w := performPostRequest(r, "/api/subscribe", requestBody)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Internal error")
	})
}

func Test_GetSubscription(t *testing.T) {
	t.Parallel()

	t.Run("successful test", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := NewMockservice(ctrl)
		server := &Server{service: mockService}

		subscriptionID := uuid.New()
		expectedSubscription := model.Subscription{ID: subscriptionID, UserID: uuid.New(), ProductID: uuid.New()}

		mockService.EXPECT().FindSubscription(gomock.Any(), subscriptionID.String()).Return(expectedSubscription, nil)

		r := gin.Default()
		r.GET("/api/subscription/:subscription_id", server.getSubscription)

		w := performRequest(r, "GET", "/api/subscription/"+subscriptionID.String())
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), subscriptionID.String())
	})

	t.Run("subscription not found", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := NewMockservice(ctrl)
		server := &Server{service: mockService}

		subscriptionID := uuid.New().String()
		mockService.EXPECT().FindSubscription(gomock.Any(), subscriptionID).Return(model.Subscription{}, fmt.Errorf("subscription not found"))

		r := gin.Default()
		r.GET("/api/subscription/:subscription_id", server.getSubscription)

		w := performRequest(r, "GET", "/api/subscription/"+subscriptionID)
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Subscription not found")
	})

	t.Run("404", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := NewMockservice(ctrl)
		server := &Server{service: mockService}

		r := gin.Default()
		r.GET("/api/subscription/:subscription_id", server.getSubscription)

		w := performRequest(r, "GET", "/api/subscription/")
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func Test_ManageSubscription(t *testing.T) {
	t.Parallel()

	t.Run("pause subscription successfully", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := NewMockservice(ctrl)
		server := &Server{service: mockService}

		subscriptionID := uuid.New().String()
		requestBody := `{"action": "pause"}`

		mockService.EXPECT().PauseSubscription(gomock.Any(), subscriptionID).Return(nil)

		r := gin.Default()
		r.POST("/api/subscription/:subscription_id/manage", server.manageSubscription)
		w := performPostRequest(r, "/api/subscription/"+subscriptionID+"/manage", requestBody)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Subscription paused")
	})

	t.Run("unpause subscription successfully", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := NewMockservice(ctrl)
		server := &Server{service: mockService}

		subscriptionID := uuid.New().String()
		requestBody := `{"action": "unpause"}`

		mockService.EXPECT().UnpauseSubscription(gomock.Any(), subscriptionID).Return(nil)

		r := gin.Default()
		r.POST("/api/subscription/:subscription_id/manage", server.manageSubscription)
		w := performPostRequest(r, "/api/subscription/"+subscriptionID+"/manage", requestBody)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Subscription unpaused")
	})

	t.Run("cancel subscription successfully", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := NewMockservice(ctrl)
		server := &Server{service: mockService}

		subscriptionID := uuid.New().String()
		requestBody := `{"action": "cancel"}`

		mockService.EXPECT().CancelSubscription(gomock.Any(), subscriptionID).Return(nil)

		r := gin.Default()
		r.POST("/api/subscription/:subscription_id/manage", server.manageSubscription)
		w := performPostRequest(r, "/api/subscription/"+subscriptionID+"/manage", requestBody)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Subscription canceled")
	})

	t.Run("invalid action", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := NewMockservice(ctrl)
		server := &Server{service: mockService}

		subscriptionID := uuid.New().String()
		requestBody := `{"action": "invalid_action"}`

		r := gin.Default()
		r.POST("/api/subscription/:subscription_id/manage", server.manageSubscription)
		w := performPostRequest(r, "/api/subscription/"+subscriptionID+"/manage", requestBody)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid action")
	})

	t.Run("missing action in request", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := NewMockservice(ctrl)
		server := &Server{service: mockService}

		subscriptionID := uuid.New().String()
		requestBody := `{}`

		r := gin.Default()
		r.POST("/api/subscription/:subscription_id/manage", server.manageSubscription)
		w := performPostRequest(r, "/api/subscription/"+subscriptionID+"/manage", requestBody)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid action")
	})

	t.Run("internal server error on pause", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := NewMockservice(ctrl)
		server := &Server{service: mockService}

		subscriptionID := uuid.New().String()
		requestBody := `{"action": "pause"}`

		mockService.EXPECT().PauseSubscription(gomock.Any(), subscriptionID).
			Return(fmt.Errorf("internal error"))

		r := gin.Default()
		r.POST("/api/subscription/:subscription_id/manage", server.manageSubscription)
		w := performPostRequest(r, "/api/subscription/"+subscriptionID+"/manage", requestBody)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "internal error")
	})
}
