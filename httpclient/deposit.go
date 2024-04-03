package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/VitoNaychev/zota-challenge/domain"
)

type HttpClient interface {
	Post(string, string, io.Reader) (*http.Response, error)
}

type DepositClient struct {
	baseURL     string
	redirectURL string
	checkoutURL string
	contentType string
	client      HttpClient
}

func NewDepositClient(baseURL, contentType string, redirectURL, checkoutURL string, client HttpClient) *DepositClient {
	return &DepositClient{
		baseURL:     baseURL,
		redirectURL: redirectURL,
		checkoutURL: checkoutURL,
		contentType: contentType,
		client:      client,
	}
}

func (d *DepositClient) Deposit(order domain.Order, customer domain.Customer) {
	depositRequest := NewDepositRequest(order, customer, d.redirectURL, fmt.Sprintf("%s?uid=%d", d.checkoutURL, order.ID))

	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(depositRequest)

	d.client.Post(d.baseURL, d.contentType, body)
}
