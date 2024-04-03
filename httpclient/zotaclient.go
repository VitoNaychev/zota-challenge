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

type ZotaClient struct {
	merchantID string
	secret     string
	endpoint   string

	baseURL     string
	contentType string

	redirectURL string
	checkoutURL string

	client HttpClient
}

func NewZotaClient(merchantID, secret, endpoint, baseURL, contentType, redirectURL, checkoutURL string, client HttpClient) *ZotaClient {
	return &ZotaClient{
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

func (z *ZotaClient) Deposit(order domain.Order, customer domain.Customer) (DepositResponseData, error) {
	signature := crypto.SignDeposit(z.endpoint, order.ID, order.Amount, customer.Email, z.secret)
	depositRequest := NewDepositRequest(order, customer, z.redirectURL, fmt.Sprintf("%s?uid=%s", z.checkoutURL, order.ID), signature)

	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(depositRequest)

	response, _ := z.client.Post(z.baseURL, z.contentType, body)

	if response.StatusCode != 200 {
		var depositErrorRespone DepositErrorResponse
		json.NewDecoder(response.Body).Decode(&depositErrorRespone)

		return DepositResponseData{}, NewDepositError(depositErrorRespone.Message)
	}

	var depositSuccessResponse DepositSuccessResponse
	json.NewDecoder(response.Body).Decode(&depositSuccessResponse)

	return depositSuccessResponse.Data, nil
}

func (z *ZotaClient) OrderStatus(orderID, merchantOrderID string) {
	params := url.Values{}

	params.Set("merchantID", z.merchantID)
	params.Set("orderID", orderID)
	params.Set("merchantOrderID", merchantOrderID)
	params.Set("timestamp", fmt.Sprint(time.Now().Unix()))
	params.Set("signature", "labadabadaba")

	urlWithParams := fmt.Sprintf("%s?%s", z.baseURL, params.Encode())
	z.client.Get(urlWithParams)
}
