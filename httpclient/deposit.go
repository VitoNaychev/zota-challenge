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
	Client HttpClient
}

func (d *DepositClient) Deposit(url, contentType, orderID string) {
	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(orderID)

	d.Client.Post(url, contentType, body)
}
