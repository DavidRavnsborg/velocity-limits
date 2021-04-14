package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

const input_src = "input.txt"
const output_src = "output-test.txt"

var fundSuccessDB = make(FundSuccessDB)
var responses = make(ResponsesTable, 0)

func main() {
	requests := loadData()

	for _, req := range requests {
		res, err := req.handleRequest()
		if err != nil {
			if err.Error() == NonUniqueIdError {
				continue
			} else {
				handleError(err)
			}
		}
		responses = append(responses, res)
	}

	var output_buffer bytes.Buffer
	for _, response := range responses {
		output_buffer.WriteString(fmt.Sprintf(`{"id":"%s","customer_id":"%s","accepted":%v}`+"\n", response.Id, response.CustomerId, response.Accepted))
	}

	ioutil.WriteFile(output_src, output_buffer.Bytes(), 0666)
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

func (req FundRequest) handleRequest() (res FundResponse, err error) {
	var reqAmount float64
	var reqTime time.Time

	reqAmount, err = strconv.ParseFloat(string(req.LoadAmount[1:]), 64)
	if err != nil {
		handleError(err)
	}

	reqTime, err = time.Parse(time.RFC3339, req.Time)
	if err != nil {
		handleError(err)
	}

	table := fundSuccessDB[req.CustomerId]
	accepted, err := checkConditions(table, responses, req.Id, req.CustomerId, reqAmount, reqTime)
	if err != nil {
		if err.Error() == NonUniqueIdError {
			return FundResponse{}, err
		} else {
			handleError(err)
		}
	}
	if accepted {
		table.updateSuccessAmount(FundSuccessRecord{
			Id:         req.Id,
			LoadAmount: reqAmount,
			Time:       reqTime,
			Accepted:   accepted,
		}, req.CustomerId)
	}

	return FundResponse{
		Id:         req.Id,
		CustomerId: req.CustomerId,
		Accepted:   accepted,
	}, nil
}
