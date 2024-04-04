package httpclient

import "github.com/VitoNaychev/zota-challenge/domain"

type ZotaClientError struct {
	msg string
}

func NewZotaClientError(msg string) *ZotaClientError {
	return &ZotaClientError{msg}
}

func (o *ZotaClientError) Error() string {
	return o.msg
}

type OrderStatusResponseExtraData struct {
	AmountChanged     bool   `json:"amountChanged"`
	AmountRounded     bool   `json:"amountRounded"`
	AmountManipulated bool   `json:"amountManipulated"`
	DCC               bool   `json:"dcc"`
	OriginalAmount    string `json:"originalAmount"`
	PaymentMethod     string `json:"paymentMethod"`
	SelectedBankCode  string `json:"selectedBankCode"`
	SelectedBankName  string `json:"selectedBankName"`
}

type OrderStatusRequest struct {
	MerchantID      string `json:"merchantID"`
	OrderID         string `json:"orderID"`
	MerchantOrderID string `json:"merchantOrderID"`
	Timestamp       string `json:"timestamp"`
}

type OrderStatusResponseData struct {
	Type                   string                       `json:"type"`
	Status                 string                       `json:"status"`
	ErrorMessage           string                       `json:"errorMessage"`
	EndpointID             string                       `json:"endpointID"`
	ProcessorTransactionID string                       `json:"processorTransactionID"`
	OrderID                string                       `json:"orderID"`
	MerchantOrderID        string                       `json:"merchantOrderID"`
	Amount                 string                       `json:"amount"`
	Currency               string                       `json:"currency"`
	CustomerEmail          string                       `json:"customerEmail"`
	CustomParam            string                       `json:"customParam"`
	ExtraData              OrderStatusResponseExtraData `json:"extraData"`
	Request                OrderStatusRequest           `json:"request"`
}

type OrderStatusSuccessResponse struct {
	Code int                     `json:"code"`
	Data OrderStatusResponseData `json:"data"`
}

type OrderStatusErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type DepositSuccessResponse struct {
	Code int                 `json:"code"`
	Data DepositResponseData `json:"data"`
}

type DepositResponseData struct {
	DepositURL      string `json:"depositURL"`
	MerchantOrderID int    `json:"merchantOrderID"`
	OrderID         string `json:"orderID"`
}

type DepositErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type DepositRequest struct {
	MerchantOrderID          string  `json:"merchantOrderID"`
	MerchantOrderDescription string  `json:"merchantOrderDescription"`
	OrderAmount              float64 `json:"orderAmount"`
	OrderCurrency            string  `json:"orderCurrency"`
	CustomerEmail            string  `json:"customerEmail"`
	CustomerFirstName        string  `json:"customerFirstName"`
	CustomerLastName         string  `json:"customerLastName"`
	CustomerAddress          string  `json:"customerAddress"`
	CustomerCountryCode      string  `json:"customerCountryCode"`
	CustomerCity             string  `json:"customerCity"`
	CustomerZipCode          string  `json:"customerZipCode"`
	CustomerPhone            string  `json:"customerPhone"`
	CustomerIP               string  `json:"customerIP"`
	RedirectURL              string  `json:"redirectURL"`
	CheckoutURL              string  `json:"checkoutURL"`
	Signature                string  `json:"signature"`
}

func NewDepositRequest(order domain.Order, customer domain.Customer, redirectURL, checkoutURL, signature string) DepositRequest {
	depositRequest := DepositRequest{
		MerchantOrderID:          order.ID,
		MerchantOrderDescription: order.Description,
		OrderAmount:              order.Amount,
		OrderCurrency:            order.Currency,
		CustomerEmail:            customer.Email,
		CustomerFirstName:        customer.FirstName,
		CustomerLastName:         customer.LastName,
		CustomerAddress:          customer.Address,
		CustomerCountryCode:      customer.CountryCode,
		CustomerCity:             customer.City,
		CustomerZipCode:          customer.ZipCode,
		CustomerPhone:            customer.Phone,
		CustomerIP:               customer.IP,
		RedirectURL:              redirectURL,
		CheckoutURL:              checkoutURL,
		Signature:                signature,
	}

	return depositRequest
}
