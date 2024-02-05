package api

import (
	"bank/utils"

	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, isOk := fl.Field().Interface().(string); isOk {
		//check whether the currency is supported
		return utils.IsCurrencySupported(currency)
	}
	return false
}
