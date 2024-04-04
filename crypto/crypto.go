package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type formatStringFunc func(format string, a ...any) string

func signDeposit(formatString formatStringFunc, endpoint, merchantOrderID string, orderAmount float64, customerEmail, secret string) string {
	s := formatString("%s%s%.2f%s%s", endpoint, merchantOrderID, orderAmount, customerEmail, secret)
	sig := sha256.Sum256([]byte(s))

	return hex.EncodeToString(sig[:])
}

func SignDeposit(endpoint, merchantOrderID string, orderAmount float64, customerEmail, secret string) string {
	return signDeposit(fmt.Sprintf, endpoint, merchantOrderID, orderAmount, customerEmail, secret)
}

func signOrderStatus(formatString formatStringFunc, merchantID, merchantOrderID string, orderID string, timestamp int64, secret string) string {
	s := formatString("%s%s%s%d%s", merchantID, merchantOrderID, orderID, timestamp, secret)
	sig := sha256.Sum256([]byte(s))

	return hex.EncodeToString(sig[:])
}

func SignOrderStatus(merchantID, merchantOrderID string, orderID string, timestamp int64, secret string) string {
	return signOrderStatus(fmt.Sprintf, merchantID, merchantOrderID, orderID, timestamp, secret)
}
