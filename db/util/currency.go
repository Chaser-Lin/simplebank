package util

const (
	USD = "USD"
	ERU = "ERU"
	RMB = "RMB"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, ERU, RMB:
		return true
	}
	return false
}
