package httpclient

import "github.com/VitoNaychev/zota-challenge/domain"

type OrderStatusError struct {
	msg string
}

func NewOrderStatusError(msg string) *OrderStatusError {
	return &OrderStatusError{msg}
}

func (o *OrderStatusError) Error() string {
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

type DepositError struct {
	msg string
}

func NewDepositError(msg string) *DepositError {
	return &DepositError{msg}
}

func (d *DepositError) Error() string {
	return d.msg
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
	MerchantOrderID          string
	MerchantOrderDescription string
	OrderAmount              float64
	OrderCurrency            string
	CustomerEmail            string
	CustomerFirstName        string
	CustomerLastName         string
	CustomerAddress          string
	CustomerCountryCode      string
	CustomerCity             string
	CustomerZipCode          string
	CustomerPhone            string
	CustomerIP               string
	RedirectURL              string
	CheckoutURL              string
	Signature                string
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
