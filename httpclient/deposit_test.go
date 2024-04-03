package httpclient_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"testing"

	"github.com/VitoNaychev/zota-challenge/crypto"
	"github.com/VitoNaychev/zota-challenge/domain"
	"github.com/VitoNaychev/zota-challenge/httpclient"
	"github.com/joho/godotenv"
)

var testOrder = domain.Order{
	ID:          "QvE8dZshpKhaOmHY",
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
	MerchantOrderID:          "QvE8dZshpKhaOmHY",
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
	RedirectURL:              "https://www.example-merchant.com/payment-return/",
	CheckoutURL:              "https://www.example-merchant.com/account/deposit/?uid=QvE8dZshpKhaOmHY",
}

var testResponseData = httpclient.DepositResponseData{
	DepositURL:      "https://api.zotapay.com/api/v1/deposit/init/8b3a6b89697e8ac8f45d964bcc90c7ba41764acd/",
	MerchantOrderID: 12,
	OrderID:         "8b3a6b89697e8ac8f45d964bcc90c7ba41764acd",
}

var testSuccessResponse = httpclient.DepositSuccessResponse{
	Code: 200,
	Data: testResponseData,
}

var testErrorResponse = httpclient.DepositErrorResponse{
	Code:    400,
	Message: "endpoint currency mismatch",
}

type StubHttpClient struct {
	url         string
	contentType string
	data        io.Reader

	code     int
	response interface{}
}

func (s *StubHttpClient) Post(url string, contentType string, data io.Reader) (*http.Response, error) {
	s.url = url
	s.contentType = contentType
	s.data = data

	response := &http.Response{}

	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(s.response)

	response.StatusCode = s.code
	response.Body = io.NopCloser(body)

	return response, nil

}

func (s *StubHttpClient) Get(url string) (*http.Response, error) {
	s.url = url
	return nil, nil
}

func TestDeposit(t *testing.T) {
	godotenv.Load("../test.env")

	merchantID := os.Getenv("MERCHANT_ID")
	secret := os.Getenv("API_SECRET_KEY")
	endpoint := os.Getenv("ENDPOINT_ID")

	baseURL := os.Getenv("BASE_URL")
	contentType := os.Getenv("CONTENT_TYPE")

	redirectURL := os.Getenv("REDIRECT_URL")
	checkoutURL := os.Getenv("CHECKOUT_URL")

	t.Run("signs request", func(t *testing.T) {
		httpClient := &StubHttpClient{
			code:     testSuccessResponse.Code,
			response: testSuccessResponse,
		}
		depositClient := httpclient.NewDepositClient(merchantID, secret, endpoint, baseURL, contentType, redirectURL, checkoutURL, httpClient)

		depositClient.Deposit(testOrder, testCustomer)

		signature := crypto.SignDeposit(endpoint, testOrder.ID, testOrder.Amount, testCustomer.Email, secret)

		var gotRequest httpclient.DepositRequest
		json.NewDecoder(httpClient.data).Decode(&gotRequest)
		AssertEqual(t, gotRequest.Signature, signature)
	})

	t.Run("sends request", func(t *testing.T) {
		httpClient := &StubHttpClient{
			code:     testSuccessResponse.Code,
			response: testSuccessResponse,
		}
		depositClient := httpclient.NewDepositClient(merchantID, secret, endpoint, baseURL, contentType, redirectURL, checkoutURL, httpClient)

		depositClient.Deposit(testOrder, testCustomer)

		AssertEqual(t, httpClient.url, baseURL)
		AssertEqual(t, httpClient.contentType, contentType)

		var gotRequest httpclient.DepositRequest
		json.NewDecoder(httpClient.data).Decode(&gotRequest)

		signature := crypto.SignDeposit(endpoint, testOrder.ID, testOrder.Amount, testCustomer.Email, secret)
		testRequest.Signature = signature

		AssertEqual(t, gotRequest, testRequest)
	})

	t.Run("returns error on failure to send request", func(t *testing.T) {})

	t.Run("returns response data on successful request", func(t *testing.T) {
		httpClient := &StubHttpClient{
			code:     testSuccessResponse.Code,
			response: testSuccessResponse,
		}
		depositClient := httpclient.NewDepositClient(merchantID, secret, endpoint, baseURL, contentType, redirectURL, checkoutURL, httpClient)

		gotResponseData, _ := depositClient.Deposit(testOrder, testCustomer)

		AssertEqual(t, gotResponseData, testResponseData)
	})

	t.Run("returns error on unsuccessful request", func(t *testing.T) {
		httpClient := &StubHttpClient{
			code:     testErrorResponse.Code,
			response: testErrorResponse,
		}
		depositClient := httpclient.NewDepositClient(merchantID, secret, endpoint, baseURL, contentType, redirectURL, checkoutURL, httpClient)

		wantErr := &httpclient.DepositError{}
		_, gotErr := depositClient.Deposit(testOrder, testCustomer)

		if !errors.As(gotErr, &wantErr) {
			t.Errorf("got error with type %v want %v", reflect.TypeOf(gotErr), reflect.TypeOf(wantErr))
		}
	})
}

func TestOrderStatus(t *testing.T) {
	godotenv.Load("../test.env")

	merchantID := os.Getenv("MERCHANT_ID")
	secret := os.Getenv("API_SECRET_KEY")
	endpoint := os.Getenv("ENDPOINT_ID")

	baseURL := os.Getenv("BASE_URL")
	contentType := os.Getenv("CONTENT_TYPE")

	redirectURL := os.Getenv("REDIRECT_URL")
	checkoutURL := os.Getenv("CHECKOUT_URL")

	t.Run("sets query parameters", func(t *testing.T) {
		merchantOrderID := "QvE8dZshpKhaOmHY"
		orderID := "8b3a6b89697e8ac8f45d964bcc90c7ba41764acd"

		httpClient := &StubHttpClient{
			code:     testSuccessResponse.Code,
			response: testSuccessResponse,
		}
		depositClient := httpclient.NewDepositClient(merchantID, secret, endpoint, baseURL, contentType, redirectURL, checkoutURL, httpClient)

		depositClient.OrderStatus(orderID, merchantOrderID)

		parsedURL, err := url.Parse(httpClient.url)
		if err != nil {
			t.Fatalf("got error %v", err)
		}

		queryParams := parsedURL.Query()

		AssertEqual(t, queryParams.Get("merchantID"), merchantID)
		AssertEqual(t, queryParams.Get("orderID"), orderID)
		AssertEqual(t, queryParams.Get("merchantOrderID"), merchantOrderID)
		AssertEqual(t, queryParams.Has("timestamp"), true)
		AssertEqual(t, queryParams.Has("signature"), true)

	})
}

func AssertEqual[T any](t testing.TB, got, want T) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
