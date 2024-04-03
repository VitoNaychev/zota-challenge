package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/VitoNaychev/zota-challenge/crypto"
	"github.com/VitoNaychev/zota-challenge/domain"
)

type HttpClient interface {
	Post(string, string, io.Reader) (*http.Response, error)
}

type DepositClient struct {
	secret   string
	endpoint string

	baseURL     string
	contentType string

	redirectURL string
	checkoutURL string

	client HttpClient
}

func NewDepositClient(secret, endpoint, baseURL, contentType, redirectURL, checkoutURL string, client HttpClient) *DepositClient {
	return &DepositClient{
		secret:   secret,
		endpoint: endpoint,

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

func (d *DepositClient) OrderStatus(orderID string, merchantOrderID int) {

}
