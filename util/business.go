package util

const (
	Withdraw = "Withdraw"
	Deposit  = "Deposit"
)

func IsSupportedBusiness(business string) bool {
	switch business {
	case Withdraw, Deposit:
		return true
	}
	return false
}
