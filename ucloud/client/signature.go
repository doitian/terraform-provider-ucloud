package client

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net/url"
	"sort"
)

// Sign signs paramters in url.Values with private key.
func GenerateSignature(params url.Values, privateKey string) string {
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	hasher := sha1.New()
	for _, k := range keys {
		values := params[k]
		for _, v := range values {
			io.WriteString(hasher, k)
			io.WriteString(hasher, v)
		}
	}

	io.WriteString(hasher, privateKey)

	return fmt.Sprintf("%x", hasher.Sum(nil))
}
