package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

const input_src = "input.txt"

// const output_src = "output.txt"
var fundSuccessDB = make(FundSuccessDB)

func main() {

	requests := loadData()

	for _, request := range requests {
		request.handleRequest()
		// fundResponse := request.handleRequest()
		// fmt.Println(fundResponse)
	}

	// TODO: Find out why this is empty (probably need to set it again somewhere)
	fmt.Println(fundSuccessDB)

	// ioutil.WriteFile("output.txt")

	// Test running a function from another file in the main package
	// config_test()
}

func loadData() []FundRequest {
	var fundRequests []FundRequest
	inputStream, err := ioutil.ReadFile(input_src)
	if err != nil {
		handleError(err)
	}

	strRequests := strings.Split(string(inputStream), "\n")
	for _, req := range strRequests {
		if len(req) == 0 {
			break
		}

		var fundRequest FundRequest
		err = json.Unmarshal([]byte(req), &fundRequest)
		if err != nil {
			handleError(err)
		}
		fundRequests = append(fundRequests, fundRequest)
	}

	return fundRequests
}

func (req FundRequest) handleRequest() (res FundResponse) {
	var reqAmount float64
	var reqTime time.Time
	var err error

	reqAmount, err = strconv.ParseFloat(string(req.LoadAmount[1:]), 64)
	if err != nil {
		handleError(err)
	}

	reqTime, err = time.Parse(time.RFC3339, req.Time)
	if err != nil {
		handleError(err)
	}

	table := fundSuccessDB[req.CustomerId]
	accepted := checkConditions(table, req.Id, reqAmount, reqTime)
	if accepted {
		fundSuccessDB[req.CustomerId] = append(table, FundSuccessRecord{
			Id:         req.Id,
			CustomerId: req.CustomerId,
			LoadAmount: reqAmount,
			Time:       reqTime,
		})
	}

	return FundResponse{
		Id:         req.Id,
		CustomerId: req.CustomerId,
		Accepted:   accepted,
	}
}
