package testdata

import (
	"github.com/VitoNaychev/zota-challenge/domain"
	"github.com/VitoNaychev/zota-challenge/httpclient"
)

var Order = domain.Order{
	ID:          "QvE8dZshpKhaOmHY",
	Description: "Test order",
	Amount:      500.00,
	Currency:    "USD",
}

var Customer = domain.Customer{
	Email:       "customer@email-address.com",
	FirstName:   "John",
	LastName:    "Doe",
	Address:     "5/5 Moo 5 Thong Nai Pan Noi Beach, Baan Tai, Koh Phangan",
	CountryCode: "TH",
	City:        "Surat Thani",
	ZipCode:     "84280",
	Phone:       "+66-77999110",
	IP:          "103.106.8.104",
}

var Request = httpclient.DepositRequest{
	MerchantOrderID:          "QvE8dZshpKhaOmHY",
	MerchantOrderDescription: "Test order",
	OrderAmount:              500.00,
	OrderCurrency:            "USD",
	CustomerEmail:            "customer@email-address.com",
	CustomerFirstName:        "John",
	CustomerLastName:         "Doe",
	CustomerAddress:          "5/5 Moo 5 Thong Nai Pan Noi Beach, Baan Tai, Koh Phangan",
	CustomerCountryCode:      "TH",
	CustomerCity:             "Surat Thani",
	CustomerZipCode:          "84280",
	CustomerPhone:            "+66-77999110",
	CustomerIP:               "103.106.8.104",
	RedirectURL:              "https://www.example-merchant.com/payment-return/",
	CheckoutURL:              "https://www.example-merchant.com/account/deposit/?uid=QvE8dZshpKhaOmHY",
}

var DepositResponseData = httpclient.DepositResponseData{
	DepositURL:      "https://api.zotapay.com/api/v1/deposit/init/8b3a6b89697e8ac8f45d964bcc90c7ba41764acd/",
	MerchantOrderID: 12,
	OrderID:         "8b3a6b89697e8ac8f45d964bcc90c7ba41764acd",
}

var DepositSuccessResponse = httpclient.DepositSuccessResponse{
	Code: 200,
	Data: DepositResponseData,
}

var DepositErrorResponse = httpclient.DepositErrorResponse{
	Code:    400,
	Message: "endpoint currency mismatch",
}

var OrderStatusRequest = httpclient.OrderStatusRequest{
	MerchantID:      "EXAMPLE-MERCHANT-ID",
	OrderID:         "8b3a6b89697e8ac8f45d964bcc90c7ba41764acd",
	MerchantOrderID: "QvE8dZshpKhaOmHY",
	Timestamp:       "1564617600",
}

var OrderStatusResponseExtraData = httpclient.OrderStatusResponseExtraData{
	AmountChanged:     true,
	AmountRounded:     true,
	AmountManipulated: false,
	DCC:               false,
	OriginalAmount:    "499.98",
	PaymentMethod:     "INSTANT-BANK-WIRE",
	SelectedBankCode:  "SCB",
	SelectedBankName:  "",
}

var OrderStatusResponseData = httpclient.OrderStatusResponseData{
	Type:                   "SALE",
	Status:                 "PROCESSING",
	ErrorMessage:           "",
	EndpointID:             "1050",
	ProcessorTransactionID: "",
	OrderID:                "8b3a6b89697e8ac8f45d964bcc90c7ba41764acd",
	MerchantOrderID:        "QvE8dZshpKhaOmHY",
	Amount:                 "500.00",
	Currency:               "THB",
	CustomerEmail:          "customer@email-address.com",
	ExtraData:              OrderStatusResponseExtraData,
	Request:                OrderStatusRequest,
}

var OrderStatusSuccessResponse = httpclient.OrderStatusSuccessResponse{
	Code: 200,
	Data: OrderStatusResponseData,
}

var OrderStatusErrorResponse = httpclient.OrderStatusErrorResponse{
	Code:    400,
	Message: "timestamp too old",
}
