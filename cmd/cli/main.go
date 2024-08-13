package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/VitoNaychev/zota-challenge/domain"
	"github.com/VitoNaychev/zota-challenge/httpclient"
	"github.com/joho/godotenv"
)

var customer = domain.Customer{
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

var order = domain.Order{
	ID:          "QvE8dZshpKhaOmHY",
	Description: "Test order",
	Amount:      500.00,
	Currency:    "USD",
}

func main() {
	godotenv.Load("../../test.env")

	config, err := httpclient.InitZotaConfigFromEnv()
	if err != nil {
		log.Fatal("InitZotaConfigFromEnv error: ", err)
	}

	zotaClient := httpclient.NewZotaClient(config, httpclient.HttpClient{})

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Input order ID: ")
		merchantOrderID, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		order.ID = merchantOrderID
		depositResponse, err := zotaClient.Deposit(order, customer)
		if err != nil {
			log.Println("Deposit error: ", err)
		}

		log.Println("Deposit call successful")
		log.Println("\tOrder ID: ", depositResponse.OrderID)
		log.Println("\tDeposit URL: ", depositResponse.DepositURL)
		log.Println("Beginning order status polling")

		for {
			orderStatusResponse, err := zotaClient.OrderStatus(depositResponse.OrderID, merchantOrderID)
			if err != nil {
				log.Println("Order Status error: ", err)
				break
			}
			if orderStatusResponse.Status == "APPROVED" ||
				orderStatusResponse.Status == "DECLINED" ||
				orderStatusResponse.Status == "ERROR" {
				log.Println("Order Status completed successfully, got status: ", orderStatusResponse.Status)
				break
			}
			log.Printf("\tCurrent status: %s; Polling...", orderStatusResponse.Status)
			time.Sleep(time.Second * 10)
		}
	}
}
