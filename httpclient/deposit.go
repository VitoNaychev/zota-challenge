package httpclient

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/VitoNaychev/zota-challenge/domain"
)

type HttpClient interface {
	Post(string, string, io.Reader) (*http.Response, error)
}

type DepositClient struct {
	url         string
	contentType string
	client      HttpClient
}

func NewDepositClient(url, contentType string, client HttpClient) *DepositClient {
	return &DepositClient{
		url:         url,
		contentType: contentType,
		client:      client,
	}
}

func (d *DepositClient) Deposit(order domain.Order, customer domain.Customer) {
	depositRequest := NewDepositRequest(order, customer)

	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(depositRequest)

	d.client.Post(d.url, d.contentType, body)
}
