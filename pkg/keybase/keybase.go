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
	keyBaseEndpoint = "https://keybase.io/_/api/1.0"
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
	endpoint := fmt.Sprintf("/user/lookup.json?key_suffix=%[1]s&fields=basics&fields=pictures", identity)
	if err := queryKeyBase(ctx, endpoint, &response); err != nil {
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
func queryKeyBase(ctx context.Context, endpoint string, data interface{}) error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", keyBaseEndpoint, endpoint), nil)
	if err != nil {
		return err
	}

	resp, err := cli.Do(ctx, req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

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
