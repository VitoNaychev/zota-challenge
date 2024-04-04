package httpclient_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	"testing"

	"github.com/VitoNaychev/zota-challenge/crypto"
	"github.com/VitoNaychev/zota-challenge/httpclient"
	"github.com/VitoNaychev/zota-challenge/httpclient/testdata"
	"github.com/joho/godotenv"
)

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

	response := &http.Response{}

	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(s.response)

	response.StatusCode = s.code
	response.Body = io.NopCloser(body)

	return response, nil
}

func TestDeposit(t *testing.T) {
	godotenv.Load("../test.env")

	config, err := httpclient.InitZotaConfigFromEnv()
	AssertNoErr(t, err)

	t.Run("signs request", func(t *testing.T) {
		httpClient := &StubHttpClient{
			code:     testdata.DepositSuccessResponse.Code,
			response: testdata.DepositSuccessResponse,
		}
		depositClient := httpclient.NewZotaClient(config, httpClient)

		depositClient.Deposit(testdata.Order, testdata.Customer)

		signature := crypto.SignDeposit(config.Endpoint, testdata.Order.ID, testdata.Order.Amount, testdata.Customer.Email, config.Secret)

		var gotRequest httpclient.DepositRequest
		json.NewDecoder(httpClient.data).Decode(&gotRequest)
		AssertEqual(t, gotRequest.Signature, signature)
	})

	t.Run("sends request", func(t *testing.T) {
		depositURLPath := "/api/v1/deposit/request/" + config.Endpoint

		httpClient := &StubHttpClient{
			code:     testdata.DepositSuccessResponse.Code,
			response: testdata.DepositSuccessResponse,
		}
		depositClient := httpclient.NewZotaClient(config, httpClient)

		depositClient.Deposit(testdata.Order, testdata.Customer)

		AssertEqual(t, httpClient.url, config.BaseURL+depositURLPath)
		AssertEqual(t, httpClient.contentType, config.ContentType)

		var gotRequest httpclient.DepositRequest
		json.NewDecoder(httpClient.data).Decode(&gotRequest)

		signature := crypto.SignDeposit(config.Endpoint, testdata.Order.ID, testdata.Order.Amount, testdata.Customer.Email, config.Secret)
		testdata.Request.Signature = signature

		AssertEqual(t, gotRequest, testdata.Request)
	})

	t.Run("returns error on failure to send request", func(t *testing.T) {})

	t.Run("returns response data on successful request", func(t *testing.T) {
		httpClient := &StubHttpClient{
			code:     testdata.DepositSuccessResponse.Code,
			response: testdata.DepositSuccessResponse,
		}
		depositClient := httpclient.NewZotaClient(config, httpClient)

		gotResponseData, _ := depositClient.Deposit(testdata.Order, testdata.Customer)

		AssertEqual(t, gotResponseData, testdata.DepositResponseData)
	})

	t.Run("returns error on unsuccessful request", func(t *testing.T) {
		httpClient := &StubHttpClient{
			code:     testdata.DepositErrorResponse.Code,
			response: testdata.DepositErrorResponse,
		}
		depositClient := httpclient.NewZotaClient(config, httpClient)

		wantErr := &httpclient.DepositError{}
		_, gotErr := depositClient.Deposit(testdata.Order, testdata.Customer)

		if !errors.As(gotErr, &wantErr) {
			t.Errorf("got error with type %v want %v", reflect.TypeOf(gotErr), reflect.TypeOf(wantErr))
		}
	})
}

func TestOrderStatus(t *testing.T) {
	godotenv.Load("../test.env")

	config, err := httpclient.InitZotaConfigFromEnv()
	AssertNoErr(t, err)

	t.Run("sets query parameters", func(t *testing.T) {
		orderStatusURLPath := "/api/v1/query/order-status/"

		merchantOrderID := "QvE8dZshpKhaOmHY"
		orderID := "8b3a6b89697e8ac8f45d964bcc90c7ba41764acd"

		httpClient := &StubHttpClient{
			code:     testdata.OrderStatusSuccessResponse.Code,
			response: testdata.OrderStatusSuccessResponse,
		}
		depositClient := httpclient.NewZotaClient(config, httpClient)

		depositClient.OrderStatus(orderID, merchantOrderID)

		parsedURL, err := url.Parse(httpClient.url)
		AssertNoErr(t, err)

		AssertEqual(t, parsedURL.Path, config.BaseURL+orderStatusURLPath)

		queryParams := parsedURL.Query()

		AssertEqual(t, queryParams.Get("merchantID"), config.MerchantID)
		AssertEqual(t, queryParams.Get("orderID"), orderID)
		AssertEqual(t, queryParams.Get("merchantOrderID"), merchantOrderID)
		AssertEqual(t, queryParams.Has("timestamp"), true)
		AssertEqual(t, queryParams.Has("signature"), true)
	})

	t.Run("signs request", func(t *testing.T) {
		merchantOrderID := "QvE8dZshpKhaOmHY"
		orderID := "8b3a6b89697e8ac8f45d964bcc90c7ba41764acd"

		httpClient := &StubHttpClient{
			code:     testdata.OrderStatusSuccessResponse.Code,
			response: testdata.OrderStatusSuccessResponse,
		}
		depositClient := httpclient.NewZotaClient(config, httpClient)

		depositClient.OrderStatus(orderID, merchantOrderID)

		parsedURL, err := url.Parse(httpClient.url)
		AssertNoErr(t, err)

		queryParams := parsedURL.Query()
		timestamp, err := strconv.ParseInt(queryParams.Get("timestamp"), 10, 64)
		AssertNoErr(t, err)

		wantSignature := crypto.SignOrderStatus(config.MerchantID, merchantOrderID, orderID, timestamp, config.Secret)
		gotSignature := queryParams.Get("signature")

		AssertEqual(t, gotSignature, wantSignature)
	})

	t.Run("returns response data on successful request", func(t *testing.T) {
		merchantOrderID := "QvE8dZshpKhaOmHY"
		orderID := "8b3a6b89697e8ac8f45d964bcc90c7ba41764acd"

		httpClient := &StubHttpClient{
			code:     testdata.OrderStatusSuccessResponse.Code,
			response: testdata.OrderStatusSuccessResponse,
		}
		depositClient := httpclient.NewZotaClient(config, httpClient)

		gotResponseData, _ := depositClient.OrderStatus(orderID, merchantOrderID)

		AssertEqual(t, gotResponseData, testdata.OrderStatusResponseData)
	})

	t.Run("returns error on unsuccessful request", func(t *testing.T) {
		merchantOrderID := "QvE8dZshpKhaOmHY"
		orderID := "8b3a6b89697e8ac8f45d964bcc90c7ba41764acd"

		httpClient := &StubHttpClient{
			code:     testdata.OrderStatusErrorResponse.Code,
			response: testdata.OrderStatusErrorResponse,
		}
		depositClient := httpclient.NewZotaClient(config, httpClient)

		wantErr := &httpclient.OrderStatusError{}
		_, gotErr := depositClient.OrderStatus(orderID, merchantOrderID)

		if !errors.As(gotErr, &wantErr) {
			t.Errorf("got error with type %v want %v", reflect.TypeOf(gotErr), reflect.TypeOf(wantErr))
		}
	})
}

func AssertNoErr(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("got error %v", err)
	}
}

func AssertEqual[T any](t testing.TB, got, want T) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
