package service

import (
	"github.com/stretchr/testify/assert"
	"gymondo/internal/model"
	"testing"
)

func Test_calculatePriceWithVoucher(t *testing.T) {
	t.Parallel()

	t.Run("successful - percentage discount", func(t *testing.T) {
		product := model.Product{
			Price:      100.0,
			Tax:        20.0,
			TotalPrice: 120.00,
		}
		voucher := model.Voucher{
			DiscountType:  model.Percentage,
			DiscountValue: 0.10,
		}

		result, err := calculatePriceWithVoucher(product, voucher)
		assert.NoError(t, err, "There should be no error")
		assert.Equal(t, 90.00, result.Price, "Price should be correctly calculated after percentage discount")
		assert.Equal(t, 18.00, result.Tax, "Tax should be correctly calculated after percentage discount")
		assert.Equal(t, 108.00, result.TotalPrice, "Total price should be the sum of price and tax")
	})

	t.Run("successful - fixed discount", func(t *testing.T) {
		product := model.Product{
			Price:      100.0,
			Tax:        20.0,
			TotalPrice: 120.00,
		}
		voucher := model.Voucher{
			DiscountType:  model.Fixed,
			DiscountValue: 10.0,
		}

		result, err := calculatePriceWithVoucher(product, voucher)
		assert.NoError(t, err, "There should be no error")
		assert.Equal(t, 91.66, result.Price, "Price should be correctly calculated after fixed discount")
		assert.Equal(t, 18.34, result.Tax, "Tax should be correctly calculated after fixed discount")
		assert.Equal(t, 110.0, result.TotalPrice, "Total price should be the sum of price and tax")
	})

	t.Run("price or tax should not be negative", func(t *testing.T) {
		product := model.Product{
			Price:      10.0,
			Tax:        5.0,
			TotalPrice: 15.00,
		}
		voucher := model.Voucher{
			DiscountType:  model.Percentage,
			DiscountValue: 1.5,
		}

		result, err := calculatePriceWithVoucher(product, voucher)
		assert.Error(t, err, "An error should be returned when price or tax becomes negative")
		assert.Equal(t, "product price couldn't be less than 0", err.Error(), "Error message should be 'product price couldn't be less 0'")
		assert.Equal(t, model.Product{}, result, "The result should be an empty product")
	})

	t.Run("price or tax should not be negative (fixed)", func(t *testing.T) {
		product := model.Product{
			Price:      10.0,
			Tax:        5.0,
			TotalPrice: 15.00,
		}
		voucher := model.Voucher{
			DiscountType:  model.Fixed,
			DiscountValue: 600.00,
		}

		result, err := calculatePriceWithVoucher(product, voucher)
		assert.Error(t, err, "An error should be returned when price or tax becomes negative")
		assert.Equal(t, "product price couldn't be less than 0", err.Error(), "Error message should be 'product price couldn't be less than 0'")
		assert.Equal(t, model.Product{}, result, "The result should be an empty product")
	})

	t.Run("zero discount", func(t *testing.T) {
		product := model.Product{
			Price:      100.0,
			Tax:        20.0,
			TotalPrice: 120.00,
		}
		voucher := model.Voucher{
			DiscountType:  model.Percentage,
			DiscountValue: 0.0,
		}

		result, err := calculatePriceWithVoucher(product, voucher)
		assert.NoError(t, err, "There should be no error")
		assert.Equal(t, product.Price, result.Price, "Price should remain the same if no discount")
		assert.Equal(t, product.Tax, result.Tax, "Tax should remain the same if no discount")
		assert.Equal(t, product.Price+product.Tax, result.TotalPrice, "Total price should be the sum of price and tax")
	})
}
