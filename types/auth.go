package types

// Account represents a chain account
type Account struct {
	Address string
	Type    string
	Height  int64
}

// NewAccount builds a new Account instance
func NewAccount(address, typeURL string, height int64) Account {
	return Account{
		Address: address,
		Type:    typeURL,
		Height:  height,
	}
}
