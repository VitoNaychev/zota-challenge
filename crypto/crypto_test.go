package crypto

import (
	"fmt"
	"testing"

	"github.com/VitoNaychev/zota-challenge/assert"
)

type spyFormatString struct {
	output string
}

func (s *spyFormatString) formatString(format string, a ...any) string {
	s.output = fmt.Sprintf(format, a...)
	return s.output
}

func TestSignDeposit(t *testing.T) {
	endpointID := "12345"
	merchantOrderID := "abcdefg"
	orderAmount := 50.00
	customerEmai := "example@customer.com"
	secretKey := "aaaa-bbbb-cccc-dddd"

	spy := spyFormatString{}
	signDeposit(spy.formatString, endpointID, merchantOrderID, orderAmount, customerEmai, secretKey)

	wantSignature := "12345abcdefg50.00example@customer.comaaaa-bbbb-cccc-dddd"
	gotSignature := spy.output

	assert.Equal(t, gotSignature, wantSignature)
}

func TestSignOrderStatus(t *testing.T) {
	merchantID := "test-merchant-id"
	merchantOrderID := "abcdefg"
	orderID := "generated-order-id"
	timestamp := int64(123456789)
	secretKey := "aaaa-bbbb-cccc-dddd"

	spy := spyFormatString{}
	signOrderStatus(spy.formatString, merchantID, merchantOrderID, orderID, timestamp, secretKey)

	wantSignature := "test-merchant-idabcdefggenerated-order-id123456789aaaa-bbbb-cccc-dddd"
	gotSignature := spy.output

	assert.Equal(t, gotSignature, wantSignature)
}
