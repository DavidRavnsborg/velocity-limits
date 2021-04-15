package main

import (
	"fmt"
	"time"
)

// Would ideally load/fetch these from another place, so we can change our limits without changing source code, unless these always remain fixed.
const dailyAmountLimit = 5000
const weeklyAmountLimit = 20000
const dailyTransactionLimit = 3

type FundSuccessDB map[string]FundSuccessTable

type FundSuccessTable []FundSuccessRecord

func (table FundSuccessTable) updateSuccessAmount(success FundSuccessRecord, customerId string) (err error) {
	fundSuccessDB[customerId] = append(table, success)
	return nil
}

func (table FundSuccessTable) isWithinDailyAmountLimit(requestTime time.Time, amountToAdd float64) bool {
	year, month, day := requestTime.Date()
	startDay := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	amount := getCumulativeAmountFunded(table, startDay, requestTime)
	amount = amount + amountToAdd
	return amount <= dailyAmountLimit
}

func (table FundSuccessTable) isWithinWeeklyAmountLimit(requestTime time.Time, amountToAdd float64) bool {
	weekDay := int(requestTime.Weekday())
	numOffset := (24 * mod(weekDay-1, 7))
	durationStr := fmt.Sprintf("-%vh", numOffset)
	startOfWeekOffset, _ := time.ParseDuration(durationStr)
	year, month, day := requestTime.Date()
	startDay := time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Add(startOfWeekOffset)
	amount := getCumulativeAmountFunded(table, startDay, requestTime)
	amount = amount + amountToAdd
	return amount <= weeklyAmountLimit
}

func (table FundSuccessTable) isWithinDailyTransactionLimit(requestTime time.Time) bool {
	year, month, day := requestTime.Date()
	startDay := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	transactions := getCumulativeTransactions(table, startDay, requestTime) + 1
	return transactions <= dailyTransactionLimit
}

func getCumulativeAmountFunded(table FundSuccessTable, startDate time.Time, endDateTime time.Time) (amount float64) {
	// startDate should start at 00:00:00 on its respective date.
	// endDateTime should be whatever time a request was sent at.
	amount = 0.0
	for _, record := range table {
		if (record.Time.After(startDate) || record.Time == startDate) && record.Time.Before(endDateTime) {
			amount += record.LoadAmount
		}
	}
	return amount
}

func getCumulativeTransactions(table FundSuccessTable, startDate time.Time, endDateTime time.Time) (transactions int) {
	// startDate should start at 00:00:00 on its respective date.
	// endTime should be whatever time a request was sent at.
	transactions = 0
	for _, record := range table {
		if (record.Time.After(startDate) || record.Time == startDate) && record.Time.Before(endDateTime) {
			transactions += 1
		}
	}
	return transactions
}
