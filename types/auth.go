package types

// Account represents a chain account
type Account struct {
	Address string
	Height  int64
}

// NewAccount builds a new Account instance
func NewAccount(address string, height int64) Account {
	return Account{
		Address: address,
		Height:  height,
	}
}
