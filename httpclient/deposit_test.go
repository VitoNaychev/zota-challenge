package httpclient_test

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/VitoNaychev/zota-challenge/httpclient"
)

type StubHttpClient struct {
	url         string
	contentType string
	data        io.Reader
}

func (s *StubHttpClient) Post(url string, contentType string, data io.Reader) (*http.Response, error) {
	s.url = url
	s.contentType = contentType
	s.data = data

	return nil, nil
}

func TestDeposit(t *testing.T) {
	t.Run("sends request to URL", func(t *testing.T) {
		httpClient := &StubHttpClient{}
		depositClient := httpclient.DepositClient{httpClient}

		url := "test-url.com"
		contentType := "application/json"
		orderID := "abcdef"

		depositClient.Deposit(url, contentType, orderID)

		AssertEqual(t, httpClient.url, url)
		AssertEqual(t, httpClient.contentType, contentType)

		var gotOrderID string
		json.NewDecoder(httpClient.data).Decode(&gotOrderID)
		AssertEqual(t, gotOrderID, orderID)

	})
}

func AssertEqual[T any](t testing.TB, got, want T) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
