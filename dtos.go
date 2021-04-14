package main

type FundRequest struct {
	Id         string `json:"id"`
	CustomerId string `json:"customer_id"`
	LoadAmount string `json:"load_amount"`
	Time       string `json:"time"`
}

type FundResponse struct {
	Id         string `json:"id"`
	CustomerId string `json:"customer_id"`
	Accepted   bool   `json:"accepted"`
}
