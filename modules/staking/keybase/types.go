package keybase

// QueryStatus contains the details of the status of a request
type QueryStatus struct {
	Name    string `json:"name"`
	ErrDesc string `json:"desc"`
	Code    int64  `json:"code"`
}

// IdentityQueryResponse represents the response to an identity query
type IdentityQueryResponse struct {
	Status  QueryStatus      `json:"status"`
	Objects []AccountDetails `json:"them"`
}

// AccountDetails contains the data of a single account details
type AccountDetails struct {
	Pictures *AccountPictures `json:"pictures"`
	ID       string           `json:"id"`
}

// AccountPictures contains the info of an account's pictures
type AccountPictures struct {
	Primary *Picture `json:"primary"`
}

// Picture contains the info of a single picture
type Picture struct {
	URL string `json:"url"`
}
