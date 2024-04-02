package httpclient

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
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

func (d *DepositClient) Deposit(orderID string) {
	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(orderID)

	d.client.Post(d.url, d.contentType, body)
}
