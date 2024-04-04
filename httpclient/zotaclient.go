package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/VitoNaychev/zota-challenge/crypto"
	"github.com/VitoNaychev/zota-challenge/domain"
)

var DEPOSIT_URL = "/api/v1/deposit/request/"
var ORDER_STATUS_URL = "/api/v1/query/order-status/"

type HttpClient interface {
	Post(string, string, io.Reader) (*http.Response, error)
	Get(string) (*http.Response, error)
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

type ZotaClient struct {
	config ZotaConfig

	client HttpClient
}

func NewZotaClient(config ZotaConfig, client HttpClient) *ZotaClient {
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

	response, _ := z.client.Post(z.config.BaseURL+DEPOSIT_URL+z.config.Endpoint, z.config.ContentType, body)

	if response.StatusCode != 200 {
		var depositErrorRespone DepositErrorResponse
		json.NewDecoder(response.Body).Decode(&depositErrorRespone)

		return DepositResponseData{}, NewDepositError(depositErrorRespone.Message)
	}

	var depositSuccessResponse DepositSuccessResponse
	json.NewDecoder(response.Body).Decode(&depositSuccessResponse)

	return depositSuccessResponse.Data, nil
}

func (z *ZotaClient) OrderStatus(orderID, merchantOrderID string) (OrderStatusResponseData, error) {
	unixTimestamp := time.Now().Unix()
	signature := crypto.SignOrderStatus(z.config.MerchantID, merchantOrderID, orderID, unixTimestamp, z.config.Secret)

	url := formatOrderStatusURL(z.config.BaseURL, z.config.MerchantID, orderID, merchantOrderID, unixTimestamp, signature)

	response, _ := z.client.Get(url)

	if response.StatusCode != 200 {
		var orderStatusErrorResponse OrderStatusErrorResponse
		json.NewDecoder(response.Body).Decode(&orderStatusErrorResponse)

		return OrderStatusResponseData{}, NewOrderStatusError(orderStatusErrorResponse.Message)
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

	urlWithParams := fmt.Sprintf("%s%s?%s", baseURL, ORDER_STATUS_URL, params.Encode())

	return urlWithParams
}
