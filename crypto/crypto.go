package crypto

import (
	"crypto/sha256"
	"fmt"
)

func SignDeposit(endpoint string, orderID int, orderAmount float64, customerEmail, secret string) string {
	s := fmt.Sprintf("%s%d%f%s%s", endpoint, orderID, orderAmount, customerEmail, secret)
	sig := sha256.Sum256([]byte(s))

	return string(sig[:])
}
