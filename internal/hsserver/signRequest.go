package hsserver

import (
	"crypto/sha256"
	"encoding/hex"
)

func signRequestV2(clientSecret, method, uri, body string) string {
	s := clientSecret + method + uri + body
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
