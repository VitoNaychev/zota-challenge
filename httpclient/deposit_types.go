package httpclient

import "github.com/VitoNaychev/zota-challenge/domain"

type DepositRequest struct {
	MerchantOrderID          int
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
}

func NewDepositRequest(order domain.Order, customer domain.Customer) DepositRequest {
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
	}

	return depositRequest
}
