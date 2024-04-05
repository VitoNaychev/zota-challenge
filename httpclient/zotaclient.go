package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/VitoNaychev/zota-challenge/crypto"
	"github.com/VitoNaychev/zota-challenge/domain"
)

const depositURL = "/api/v1/deposit/request/"
const orderStatusURL = "/api/v1/query/order-status/"

type ZotaConfigError struct {
	msg string
}

func NewZotaConfigError(name string) *ZotaConfigError {
	return &ZotaConfigError{fmt.Sprintf("enviornment variable %s is not set", name)}
}

func (z *ZotaConfigError) Error() string {
	return z.msg
}

type ZotaConfig struct {
	MerchantID string
	Secret     string
	Endpoint   string

	BaseURL     string
	ContentType string

	RedirectURL string
	CheckoutURL string
}

func InitZotaConfigFromEnv() (ZotaConfig, error) {
	if err := requireEnvVariable("MERCHANT_ID"); err != nil {
		return ZotaConfig{}, err
	}
	if err := requireEnvVariable("API_SECRET_KEY"); err != nil {
		return ZotaConfig{}, err
	}
	if err := requireEnvVariable("ENDPOINT_ID"); err != nil {
		return ZotaConfig{}, err
	}

	if err := requireEnvVariable("BASE_URL"); err != nil {
		return ZotaConfig{}, err
	}
	if err := requireEnvVariable("CONTENT_TYPE"); err != nil {
		return ZotaConfig{}, err
	}

	if err := requireEnvVariable("REDIRECT_URL"); err != nil {
		return ZotaConfig{}, err
	}
	if err := requireEnvVariable("CHECKOUT_URL"); err != nil {
		return ZotaConfig{}, err
	}

	config := ZotaConfig{
		MerchantID: os.Getenv("MERCHANT_ID"),
		Secret:     os.Getenv("API_SECRET_KEY"),
		Endpoint:   os.Getenv("ENDPOINT_ID"),

		BaseURL:     os.Getenv("BASE_URL"),
		ContentType: os.Getenv("CONTENT_TYPE"),

		RedirectURL: os.Getenv("REDIRECT_URL"),
		CheckoutURL: os.Getenv("CHECKOUT_URL"),
	}

	return config, nil
}

func requireEnvVariable(name string) error {
	if _, ok := os.LookupEnv(name); !ok {
		return NewZotaConfigError(name)
	}
	return nil
}

type Client interface {
	Post(string, string, io.Reader) (*http.Response, error)
	Get(string) (*http.Response, error)
}

type HttpClient struct{}

func (h HttpClient) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	return http.Post(url, contentType, body)
}

func (h HttpClient) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

type ZotaClient struct {
	config ZotaConfig

	client Client
}

func NewZotaClient(config ZotaConfig, client Client) *ZotaClient {
	return &ZotaClient{
		config: config,
		client: client,
	}
}

func (z *ZotaClient) Deposit(order domain.Order, customer domain.Customer) (DepositResponseData, error) {
	signature := crypto.SignDeposit(z.config.Endpoint, order.ID, order.Amount, customer.Email, z.config.Secret)
	depositRequest := NewDepositRequest(order, customer, z.config.RedirectURL, fmt.Sprintf("%s?uid=%s", z.config.CheckoutURL, order.ID), signature)

	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(depositRequest)

	response, err := z.client.Post(z.config.BaseURL+depositURL+z.config.Endpoint, z.config.ContentType, body)
	if err != nil {
		return DepositResponseData{}, err
	}

	if response.StatusCode != 200 {
		var depositErrorRespone DepositErrorResponse
		json.NewDecoder(response.Body).Decode(&depositErrorRespone)

		return DepositResponseData{}, NewZotaClientError(depositErrorRespone.Message)
	}

	var depositSuccessResponse DepositSuccessResponse
	json.NewDecoder(response.Body).Decode(&depositSuccessResponse)

	return depositSuccessResponse.Data, nil
}

func (z *ZotaClient) OrderStatus(orderID, merchantOrderID string) (OrderStatusResponseData, error) {
	unixTimestamp := time.Now().Unix()
	signature := crypto.SignOrderStatus(z.config.MerchantID, merchantOrderID, orderID, unixTimestamp, z.config.Secret)

	url := formatOrderStatusURL(z.config.BaseURL, z.config.MerchantID, orderID, merchantOrderID, unixTimestamp, signature)

	response, err := z.client.Get(url)
	if err != nil {
		return OrderStatusResponseData{}, err
	}

	if response.StatusCode != 200 {
		var orderStatusErrorResponse OrderStatusErrorResponse
		json.NewDecoder(response.Body).Decode(&orderStatusErrorResponse)

		return OrderStatusResponseData{}, NewZotaClientError(orderStatusErrorResponse.Message)
	}

	var orderStatusSuccessResponse OrderStatusSuccessResponse
	json.NewDecoder(response.Body).Decode(&orderStatusSuccessResponse)

	return orderStatusSuccessResponse.Data, nil
}

func formatOrderStatusURL(baseURL, merchantID, orderID, merchantOrderID string, timestamp int64, signature string) string {
	params := url.Values{}
	params.Set("merchantID", merchantID)
	params.Set("orderID", orderID)
	params.Set("merchantOrderID", merchantOrderID)
	params.Set("timestamp", fmt.Sprint(timestamp))
	params.Set("signature", signature)

	urlWithParams := fmt.Sprintf("%s%s?%s", baseURL, orderStatusURL, params.Encode())

	return urlWithParams
}
