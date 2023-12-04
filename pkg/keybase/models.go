package keybase

type (
	// IdentityQueryResponse represents the response to an identity query
	IdentityQueryResponse struct {
		Status  QueryStatus      `json:"status"`
		Objects []AccountDetails `json:"them"`
	}

	// QueryStatus contains the details of the status of a request
	QueryStatus struct {
		ErrDesc string `json:"desc"`
		Code    int64  `json:"code"`
	}

	// AccountDetails contains the data of a single account details
	AccountDetails struct {
		Pictures *AccountPictures `json:"pictures"`
	}

	// AccountPictures contains the info of an account's pictures
	AccountPictures struct {
		Primary *Picture `json:"primary"`
	}

	// Picture contains the info of a single picture
	Picture struct {
		URL string `json:"url"`
	}
)
