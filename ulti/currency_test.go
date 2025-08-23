package ulti

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsSupportedCurrency(t *testing.T) {
	validCurrency := RandomCurrency()
	invalidCurrency := "ABC"

	isValid := IsSupportedCurrency(validCurrency)
	require.Equal(t, true, isValid)

	isValid = IsSupportedCurrency(invalidCurrency)
	require.Equal(t, false, isValid)
}
