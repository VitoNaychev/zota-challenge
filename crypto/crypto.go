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

func SignOrderStatus(merchantID, merchantOrderID string, orderID string, timestamp int64, secret string) string {
	s := fmt.Sprintf("%s%s%s%d%s", merchantID, merchantOrderID, orderID, timestamp, secret)
	sig := sha256.Sum256([]byte(s))

	return hex.EncodeToString(sig[:])
}
