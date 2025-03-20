package service

import (
	"errors"
	"gymondo/internal/model"
	"math"
)

func calculatePriceWithVoucher(
	product model.Product,
	voucher model.Voucher,
) (model.Product, error) {
	switch voucher.DiscountType {
	case model.Percentage:
		product.Price -= product.Price * voucher.DiscountValue
		product.Tax -= product.Tax * voucher.DiscountValue
	case model.Fixed:
		discountPercent := voucher.DiscountValue / product.TotalPrice
		product.Price -= product.Price * discountPercent
		product.Tax -= product.Tax * discountPercent
	}

	if product.Price < 0 {
		return model.Product{}, errors.New("product price couldn't be less than 0")
	}
	if product.Tax < 0 {
		return model.Product{}, errors.New("tax couldn't be less than 0")
	}

	product.Price = math.Floor(product.Price*100) / 100
	product.Tax = math.Ceil(product.Tax*100) / 100
	product.TotalPrice = product.Price + product.Tax

	return product, nil
}
