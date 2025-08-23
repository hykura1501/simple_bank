package ulti

import (
	"slices"
)

const (
	USD = "USD"
	VND = "VND"
	EUR = "EUR"
	CAD = "CAD"
)

var CURRENCIES = []string{
	CAD,
	VND,
	EUR,
	USD,
}

func IsSupportedCurrency(currency string) bool {
	return slices.Contains(CURRENCIES, currency)
}
