package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/VitoNaychev/zota-challenge/domain"
	"github.com/VitoNaychev/zota-challenge/httpclient"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../../.env")

	config, err := httpclient.InitZotaConfigFromEnv()
	if err != nil {
		log.Fatal("InitZotaConfigFromEnv error: ", err)
	}

	zotaClient := httpclient.NewZotaClient(config, httpclient.HttpClient{})

	var customer domain.Customer
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println()
		fmt.Println("[1] Create/Update customer    [2] Create deposit")
		fmt.Print("Choose and option: ")

		optionStr, _ := reader.ReadString('\n')
		option, err := strconv.Atoi(optionStr)
		if err != nil {
			fmt.Println("Invalid option")
			continue
		}

		switch option {
		case 1:
			customer = createCustomer()
		case 2:
			order := createOrder()

			zotaClient.Deposit(order, customer)
		}

	}

}

func createCustomer() domain.Customer {
	customer := domain.Customer{
		Address:     "Sofia Center, General Gurko St 21, 1000 Sofia",
		CountryCode: "BG",
		City:        "Sofia",
		ZipCode:     "1000",
		Phone:       "0893 885 158",
		IP:          "192.168.0.10",
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\tEnter Email: ")
	customer.Email, _ = reader.ReadString('\n')
	fmt.Print("\tEnter First Name: ")
	customer.FirstName, _ = reader.ReadString('\n')
	fmt.Print("\tEnter Last Name: ")
	customer.LastName, _ = reader.ReadString('\n')

	return customer
}

func createOrder() domain.Order {
	order := domain.Order{}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\tEnter Order ID: ")
	order.ID, _ = reader.ReadString('\n')
	fmt.Print("\tEnter Order Description: ")
	order.Description, _ = reader.ReadString('\n')
	fmt.Print("\tEnter Order Amount: ")
	amountStr, _ := reader.ReadString('\n')
	order.Amount, _ = strconv.ParseFloat(amountStr, 64)

	order.Currency = os.Getenv("CURRENCY")

	return order
}
