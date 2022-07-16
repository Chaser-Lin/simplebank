package api

import (
	util2 "SimpleBank/util"
	"github.com/go-playground/validator/v10"
)

var ValidCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return util2.IsSupportedCurrency(currency)
	}
	return false
}

var ValidBusiness validator.Func = func(fl validator.FieldLevel) bool {
	if business, ok := fl.Field().Interface().(string); ok {
		return util2.IsSupportedBusiness(business)
	}
	return false
}
