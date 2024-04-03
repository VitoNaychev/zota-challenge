package httpclient_test

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/VitoNaychev/zota-challenge/domain"
	"github.com/VitoNaychev/zota-challenge/httpclient"
)

var testOrder = domain.Order{
	ID:          12,
	Description: "Test order",
	Amount:      500.00,
	Currency:    "USD",
}

var testCustomer = domain.Customer{
	Email:       "customer@email-address.com",
	FirstName:   "John",
	LastName:    "Doe",
	Address:     "5/5 Moo 5 Thong Nai Pan Noi Beach, Baan Tai, Koh Phangan",
	CountryCode: "TH",
	City:        "Surat Thani",
	ZipCode:     "84280",
	Phone:       "+66-77999110",
	IP:          "103.106.8.104",
}

var testRequest = httpclient.DepositRequest{
	MerchantOrderID:          12,
	MerchantOrderDescription: "Test order",
	OrderAmount:              500.00,
	OrderCurrency:            "USD",
	CustomerEmail:            "customer@email-address.com",
	CustomerFirstName:        "John",
	CustomerLastName:         "Doe",
	CustomerAddress:          "5/5 Moo 5 Thong Nai Pan Noi Beach, Baan Tai, Koh Phangan",
	CustomerCountryCode:      "TH",
	CustomerCity:             "Surat Thani",
	CustomerZipCode:          "84280",
	CustomerPhone:            "+66-77999110",
	CustomerIP:               "103.106.8.104",
}

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
		url := "test-url.com"
		contentType := "application/json"

		httpClient := &StubHttpClient{}
		depositClient := httpclient.NewDepositClient(url, contentType, httpClient)

		depositClient.Deposit(testOrder, testCustomer)

		AssertEqual(t, httpClient.url, url)
		AssertEqual(t, httpClient.contentType, contentType)

		var gotRequest httpclient.DepositRequest
		json.NewDecoder(httpClient.data).Decode(&gotRequest)
		AssertEqual(t, gotRequest, testRequest)
	})
}

func AssertEqual[T any](t testing.TB, got, want T) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
