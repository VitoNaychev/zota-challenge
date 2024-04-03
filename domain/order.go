package domain

type Order struct {
	ID          int
	Description string
	Amount      float64
	Currency    string
}
