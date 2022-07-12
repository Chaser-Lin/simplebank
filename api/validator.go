package api

import (
	"SimpleBank/db/util"
	"github.com/go-playground/validator/v10"
)

var ValidCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}

var ValidBusiness validator.Func = func(fl validator.FieldLevel) bool {
	if business, ok := fl.Field().Interface().(string); ok {
		return util.IsSupportedBusiness(business)
	}
	return false
}
