package hsserver

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func signRequestV2(clientSecret, method, uri string, body []byte) string {
	s := clientSecret + method + uri + string(body)
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func signRequestV3(clientSecret, method, uri string, body []byte, timestamp int64) string {
	str := fmt.Sprintf("%s%s%s%d", method, uri, string(body), timestamp)
	h := hmac.New(sha256.New, []byte(clientSecret))
	h.Write([]byte(str))
	encodedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return encodedStr
}
