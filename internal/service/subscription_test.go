package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gymondo/internal/model"
	"testing"
	"time"
)

func Test_Service_Subscribe(t *testing.T) {
	t.Parallel()

	t.Run("failed to fetch user", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		userID := uuid.New().String()
		productID := uuid.New().String()

		mockRepo.EXPECT().GetUser(gomock.Any(), userID).Return(model.User{}, fmt.Errorf("database error"))

		expectedError := "failed to fetch user"
		_, err := service.Subscribe(context.Background(), userID, productID, "", false)
		assert.Errorf(t, err, expectedError)
	})

	t.Run("failed to fetch product", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		userID := uuid.New()
		productID := uuid.New().String()

		mockRepo.EXPECT().GetUser(gomock.Any(), userID.String()).Return(model.User{ID: userID}, nil)
		mockRepo.EXPECT().GetProduct(gomock.Any(), productID).Return(model.Product{}, fmt.Errorf("database error"))

		expectedError := "failed to fetch product"
		_, err := service.Subscribe(context.Background(), userID.String(), productID, "", false)
		assert.Errorf(t, err, expectedError)
	})

	t.Run("failed to fetch voucher", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		userID := uuid.New()
		productID := uuid.New()
		voucherCode := "voucher123"

		mockRepo.EXPECT().GetUser(gomock.Any(), userID.String()).Return(model.User{ID: userID}, nil)
		mockRepo.EXPECT().GetProduct(gomock.Any(), productID.String()).Return(model.Product{ID: productID, DurationDays: 30, Price: 100, Tax: 10, TotalPrice: 110}, nil)
		mockRepo.EXPECT().GetVoucherByCode(gomock.Any(), voucherCode).Return(model.Voucher{}, fmt.Errorf("voucher not found"))

		expectedError := "failed to fetch voucher"
		_, err := service.Subscribe(context.Background(), userID.String(), productID.String(), voucherCode, false)
		assert.Errorf(t, err, expectedError)
	})

	t.Run("successful subscription without voucher", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		userID := uuid.New()
		productID := uuid.New()

		mockRepo.EXPECT().GetUser(gomock.Any(), userID.String()).Return(model.User{ID: userID}, nil)
		mockRepo.EXPECT().GetProduct(gomock.Any(), productID.String()).Return(model.Product{ID: productID, DurationDays: 30, Price: 100, Tax: 10, TotalPrice: 110}, nil)
		mockRepo.EXPECT().SaveSubscription(gomock.Any(), gomock.Any()).Return(nil)

		subscriptionID, err := service.Subscribe(context.Background(), userID.String(), productID.String(), "", false)
		assert.NoError(t, err)
		assert.NotEmpty(t, subscriptionID)
	})

	t.Run("successful subscription with voucher", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		userID := uuid.New()
		productID := uuid.New()
		voucherCode := "voucher123"

		mockRepo.EXPECT().GetUser(gomock.Any(), userID.String()).Return(model.User{ID: userID}, nil)
		mockRepo.EXPECT().GetProduct(gomock.Any(), productID.String()).Return(model.Product{ID: productID, DurationDays: 30, Price: 100, Tax: 10, TotalPrice: 110}, nil)
		mockRepo.EXPECT().GetVoucherByCode(gomock.Any(), voucherCode).Return(model.Voucher{DiscountType: model.Fixed, DiscountValue: 10}, nil)
		mockRepo.EXPECT().SaveSubscription(gomock.Any(), gomock.Any()).Return(nil)

		subscriptionID, err := service.Subscribe(context.Background(), userID.String(), productID.String(), voucherCode, false)
		assert.NoError(t, err)
		assert.NotEmpty(t, subscriptionID)
	})

	t.Run("successful subscription with trial period", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		userID := uuid.New()
		productID := uuid.New()

		mockRepo.EXPECT().GetUser(gomock.Any(), userID.String()).Return(model.User{ID: userID}, nil)
		mockRepo.EXPECT().GetProduct(gomock.Any(), productID.String()).Return(model.Product{ID: productID, DurationDays: 30, Price: 100, Tax: 10, TotalPrice: 110}, nil)
		mockRepo.EXPECT().SaveSubscription(gomock.Any(), gomock.Any()).Return(nil)

		subscriptionID, err := service.Subscribe(context.Background(), userID.String(), productID.String(), "", true)
		assert.NoError(t, err)
		assert.NotEmpty(t, subscriptionID)
	})
}

func Test_Service_FindSubscription(t *testing.T) {
	t.Parallel()

	t.Run("failed to find subscription", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		subscriptionID := "sub123"
		mockRepo.EXPECT().GetSubscription(gomock.Any(), subscriptionID).Return(model.Subscription{}, fmt.Errorf("database error"))

		_, err := service.FindSubscription(context.Background(), subscriptionID)
		expectedError := fmt.Sprintf("failed to fetch subscription with ID %s: database error", subscriptionID)
		assert.EqualError(t, err, expectedError)
	})

	t.Run("successful fetch subscription", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		subscriptionID := uuid.New()
		expectedSubscription := model.Subscription{
			ID:     subscriptionID,
			Status: model.Active,
		}
		mockRepo.EXPECT().GetSubscription(gomock.Any(), subscriptionID.String()).Return(expectedSubscription, nil)

		subscription, err := service.FindSubscription(context.Background(), subscriptionID.String())
		assert.NoError(t, err)
		assert.Equal(t, expectedSubscription, subscription)
	})
}

func Test_Service_PauseSubscription(t *testing.T) {
	t.Parallel()

	t.Run("failed to find subscription for pause", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		subscriptionID := uuid.New()
		mockRepo.EXPECT().GetSubscription(gomock.Any(), subscriptionID.String()).Return(model.Subscription{}, fmt.Errorf("database error"))

		expectedError := "failed to find subscription"
		err := service.PauseSubscription(context.Background(), subscriptionID.String())
		assert.Errorf(t, err, expectedError)
	})

	t.Run("subscription already paused", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		subscriptionID := uuid.New()
		subscription := model.Subscription{
			ID:     subscriptionID,
			Status: model.Paused,
		}
		mockRepo.EXPECT().GetSubscription(gomock.Any(), subscriptionID.String()).Return(subscription, nil)

		expectedError := "subscription is already paused"
		err := service.PauseSubscription(context.Background(), subscriptionID.String())
		assert.EqualError(t, err, expectedError)
	})

	t.Run("subscription is cancelled", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		subscriptionID := uuid.New()
		subscription := model.Subscription{
			ID:     subscriptionID,
			Status: model.Canceled,
		}
		mockRepo.EXPECT().GetSubscription(gomock.Any(), subscriptionID.String()).Return(subscription, nil)

		expectedError := "subscription is canceled"
		err := service.PauseSubscription(context.Background(), subscriptionID.String())
		assert.EqualError(t, err, expectedError)
	})

	t.Run("successful pause subscription", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		subscriptionID := uuid.New()
		subscription := model.Subscription{
			ID:     subscriptionID,
			Status: model.Active,
		}

		mockRepo.EXPECT().GetSubscription(gomock.Any(), subscriptionID.String()).Return(subscription, nil)
		mockRepo.EXPECT().UpdateSubscription(gomock.Any(), gomock.Any()).Return(nil)

		err := service.PauseSubscription(context.Background(), subscriptionID.String())
		assert.NoError(t, err)
	})

	t.Run("fail update subscription", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		subscriptionID := uuid.New()
		subscription := model.Subscription{
			ID:     subscriptionID,
			Status: model.Active,
		}

		mockRepo.EXPECT().GetSubscription(gomock.Any(), subscriptionID.String()).Return(subscription, nil)

		expectedError := errors.New("test error")
		mockRepo.EXPECT().UpdateSubscription(gomock.Any(), gomock.Any()).Return(expectedError)

		err := service.PauseSubscription(context.Background(), subscriptionID.String())
		assert.ErrorIs(t, err, expectedError)
	})

	t.Run("subscription in trial period", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		subscriptionID := uuid.New()
		trialEndDate := time.Now().Add(24 * time.Hour)
		subscription := model.Subscription{
			ID:           subscriptionID,
			Status:       model.Active,
			TrialEndDate: &trialEndDate,
		}
		mockRepo.EXPECT().GetSubscription(gomock.Any(), subscriptionID.String()).Return(subscription, nil)

		expectedError := errors.New("can't pause subscription during trial period")
		err := service.PauseSubscription(context.Background(), subscriptionID.String())
		assert.EqualError(t, err, expectedError.Error())
	})
}

func Test_Service_UnpauseSubscription(t *testing.T) {
	t.Parallel()

	t.Run("failed to find subscription for unpause", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		subscriptionID := uuid.New()
		mockRepo.EXPECT().GetSubscription(gomock.Any(), subscriptionID.String()).Return(model.Subscription{}, fmt.Errorf("database error"))

		expectedError := "failed to find subscription"
		err := service.UnpauseSubscription(context.Background(), subscriptionID.String())
		assert.Errorf(t, err, expectedError)
	})

	t.Run("subscription is already active", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		subscriptionID := uuid.New()
		subscription := model.Subscription{
			ID:     subscriptionID,
			Status: model.Active,
		}
		mockRepo.EXPECT().GetSubscription(gomock.Any(), subscriptionID.String()).Return(subscription, nil)

		expectedError := "subscription is already active"
		err := service.UnpauseSubscription(context.Background(), subscriptionID.String())
		assert.EqualError(t, err, expectedError)
	})

	t.Run("subscription is canceled", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		subscriptionID := uuid.New()
		subscription := model.Subscription{
			ID:     subscriptionID,
			Status: model.Canceled,
		}
		mockRepo.EXPECT().GetSubscription(gomock.Any(), subscriptionID.String()).Return(subscription, nil)

		expectedError := "subscription is canceled"
		err := service.UnpauseSubscription(context.Background(), subscriptionID.String())
		assert.EqualError(t, err, expectedError)
	})

	t.Run("subscription is paused", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		subscriptionID := uuid.New()
		subscription := model.Subscription{
			ID:     subscriptionID,
			Status: model.Paused,
		}

		mockRepo.EXPECT().GetSubscription(gomock.Any(), subscriptionID.String()).Return(subscription, nil)
		mockRepo.EXPECT().UpdateSubscription(gomock.Any(), gomock.Any()).Return(nil)

		err := service.UnpauseSubscription(context.Background(), subscriptionID.String())
		assert.NoError(t, err)
	})

	t.Run("fail update subscription", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		subscriptionID := uuid.New()
		subscription := model.Subscription{
			ID:     subscriptionID,
			Status: model.Paused,
		}

		mockRepo.EXPECT().GetSubscription(gomock.Any(), subscriptionID.String()).Return(subscription, nil)

		expectedError := errors.New("test error")
		mockRepo.EXPECT().UpdateSubscription(gomock.Any(), gomock.Any()).Return(expectedError)

		err := service.UnpauseSubscription(context.Background(), subscriptionID.String())
		assert.ErrorIs(t, err, expectedError)
	})
}

func Test_Service_CancelSubscription(t *testing.T) {
	t.Parallel()

	t.Run("failed to find subscription for cancel", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		subscriptionID := uuid.New()
		mockRepo.EXPECT().GetSubscription(gomock.Any(), subscriptionID.String()).Return(model.Subscription{}, fmt.Errorf("database error"))

		expectedError := "failed to find subscription"
		err := service.CancelSubscription(context.Background(), subscriptionID.String())
		assert.Errorf(t, err, expectedError)
	})

	t.Run("subscription is paused", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		subscriptionID := uuid.New()
		subscription := model.Subscription{
			ID:     subscriptionID,
			Status: model.Paused,
		}
		mockRepo.EXPECT().GetSubscription(gomock.Any(), subscriptionID.String()).Return(subscription, nil)

		expectedError := "subscription is paused"
		err := service.CancelSubscription(context.Background(), subscriptionID.String())
		assert.EqualError(t, err, expectedError)
	})

	t.Run("subscription is already canceled", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		subscriptionID := uuid.New()
		subscription := model.Subscription{
			ID:     subscriptionID,
			Status: model.Canceled,
		}
		mockRepo.EXPECT().GetSubscription(gomock.Any(), subscriptionID.String()).Return(subscription, nil)

		expectedError := "subscription is already canceled"
		err := service.CancelSubscription(context.Background(), subscriptionID.String())
		assert.EqualError(t, err, expectedError)
	})

	t.Run("successful cancel subscription", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		subscriptionID := uuid.New()
		subscription := model.Subscription{
			ID:     subscriptionID,
			Status: model.Active,
		}

		mockRepo.EXPECT().GetSubscription(gomock.Any(), subscriptionID.String()).Return(subscription, nil)
		mockRepo.EXPECT().UpdateSubscription(gomock.Any(), gomock.Any()).Return(nil)

		err := service.CancelSubscription(context.Background(), subscriptionID.String())
		assert.NoError(t, err)
	})

	t.Run("fail update subscription", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepository(ctrl)
		service := &Service{repository: mockRepo}

		subscriptionID := uuid.New()
		subscription := model.Subscription{
			ID:     subscriptionID,
			Status: model.Active,
		}

		mockRepo.EXPECT().GetSubscription(gomock.Any(), subscriptionID.String()).Return(subscription, nil)

		expectedError := errors.New("test error")
		mockRepo.EXPECT().UpdateSubscription(gomock.Any(), gomock.Any()).Return(expectedError)

		err := service.CancelSubscription(context.Background(), subscriptionID.String())
		assert.ErrorIs(t, err, expectedError)
	})
}
