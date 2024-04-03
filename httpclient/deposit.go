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

type HttpClient interface {
	Post(string, string, io.Reader) (*http.Response, error)
	Get(string) (*http.Response, error)
}

type DepositClient struct {
	merchantID string
	secret     string
	endpoint   string

	baseURL     string
	contentType string

	redirectURL string
	checkoutURL string

	client HttpClient
}

func NewDepositClient(merchantID, secret, endpoint, baseURL, contentType, redirectURL, checkoutURL string, client HttpClient) *DepositClient {
	return &DepositClient{
		merchantID: merchantID,
		secret:     secret,
		endpoint:   endpoint,

		baseURL:     baseURL,
		contentType: contentType,

		redirectURL: redirectURL,
		checkoutURL: checkoutURL,

		client: client,
	}
}

func (d *DepositClient) Deposit(order domain.Order, customer domain.Customer) (DepositResponseData, error) {
	signature := crypto.SignDeposit(d.endpoint, order.ID, order.Amount, customer.Email, d.secret)
	depositRequest := NewDepositRequest(order, customer, d.redirectURL, fmt.Sprintf("%s?uid=%s", d.checkoutURL, order.ID), signature)

	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(depositRequest)

	response, _ := d.client.Post(d.baseURL, d.contentType, body)

	if response.StatusCode != 200 {
		var depositErrorRespone DepositErrorResponse
		json.NewDecoder(response.Body).Decode(&depositErrorRespone)

		return DepositResponseData{}, NewDepositError(depositErrorRespone.Message)
	}

	var depositSuccessResponse DepositSuccessResponse
	json.NewDecoder(response.Body).Decode(&depositSuccessResponse)

	return depositSuccessResponse.Data, nil
}

func (d *DepositClient) OrderStatus(orderID, merchantOrderID string) {
	params := url.Values{}

	params.Set("merchantID", d.merchantID)
	params.Set("orderID", orderID)
	params.Set("merchantOrderID", merchantOrderID)
	params.Set("timestamp", fmt.Sprint(time.Now().Unix()))
	params.Set("signature", "labadabadaba")

	urlWithParams := fmt.Sprintf("%s?%s", d.baseURL, params.Encode())
	d.client.Get(urlWithParams)
}
