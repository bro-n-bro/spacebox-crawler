package keybase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

const (
	fmtKeyBaseEndpoint = "https://keybase.io/_/api/1.0/user/lookup.json?key_suffix=%[1]s&fields=basics&fields=pictures"
)

var (
	errInvalidStatusCode = errors.New("invalid status code")
)

// GetAvatarURL returns the avatar URL from the given identity.
// If no identity is found, it returns an empty string instead.
func GetAvatarURL(ctx context.Context, identity string) (string, error) {
	if len(identity) != 16 {
		return "", nil
	}

	var response IdentityQueryResponse
	if err := queryKeyBase(ctx, identity, &response); err != nil {
		return "", err
	}

	// The server responded with an error
	if response.Status.Code != 0 {
		return "", fmt.Errorf("%w: %s", errInvalidStatusCode, response.Status.ErrDesc)
	}

	// No images found
	if len(response.Objects) == 0 {
		return "", nil
	}

	// Either the pictures do not exist, or the primary one does not exist, or the URL is empty
	data := response.Objects[0]
	if data.Pictures == nil || data.Pictures.Primary == nil || len(data.Pictures.Primary.URL) == 0 {
		return "", nil
	}

	// The picture URL is found
	return data.Pictures.Primary.URL, nil
}

// queryKeyBase queries the KeyBase APIs for the given endpoint, and deserializes
// the response as a JSON object inside the given data.
// Uses custom HTTP client to rate limit the requests to avoid 429 error code from API.
func queryKeyBase(ctx context.Context, identity string, data interface{}) error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(fmtKeyBaseEndpoint, identity), nil)
	if err != nil {
		return err
	}

	// call the API
	resp, err := DefaultHTTPClient.Do(ctx, req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// check status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: %v", errInvalidStatusCode, resp.Status)
	}

	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err = jsoniter.Unmarshal(bz, &data); err != nil {
		return err
	}

	return nil
}
