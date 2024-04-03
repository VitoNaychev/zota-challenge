package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func SignDeposit(endpoint, merchantOrderID string, orderAmount float64, customerEmail, secret string) string {
	s := fmt.Sprintf("%s%s%f%s%s", endpoint, merchantOrderID, orderAmount, customerEmail, secret)
	sig := sha256.Sum256([]byte(s))

	return hex.EncodeToString(sig[:])
}
