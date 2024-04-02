package handler_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/VitoNaychev/zota-challenge/handler"
)

func TestPaymentHandler(t *testing.T) {
	t.Run("returns Accepted on POST request", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/pay/", nil)
		response := httptest.NewRecorder()

		handler.PaymentHandler(response, request)

		AssertEqual(t, response.Code, http.StatusAccepted)
	})
}

func AssertEqual[T any](t testing.TB, got, want T) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
