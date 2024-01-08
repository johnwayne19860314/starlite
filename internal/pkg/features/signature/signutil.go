package signature

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"strings"
)

// Sign for input info with sigSecret+HMAC-MD5, Output is upper case
func Sign(sigSecret, sigData string) (string, error) {
	if sigSecret == "" || sigData == "" {
		return "", errors.New("both sigSecret and sigData must be non-empty")
	}

	h := hmac.New(md5.New, []byte(sigSecret))
	if _, err := h.Write([]byte(sigData)); err != nil {
		return "", err
	}
	signatureBytes := h.Sum(nil)
	return strings.ToUpper(hex.EncodeToString(signatureBytes)), nil
}
