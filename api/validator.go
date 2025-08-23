package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/hykura1501/simple_bank/ulti"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	currency := fl.Field().String()
	return ulti.IsSupportedCurrency(currency)
}
