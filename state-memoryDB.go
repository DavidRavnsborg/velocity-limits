package main

import "time"

type FundSuccessDB map[string]FundSuccessTable

type FundSuccessTable []FundSuccessRecord

type FundSuccessRecord struct {
	Id         string
	CustomerId string
	LoadAmount float64
	Time       time.Time
}

func (table FundSuccessTable) isUnderDailyAmount() bool {
	return true
}

func (table FundSuccessTable) isUnderWeeklyAmount() bool {
	return true
}

func (table FundSuccessTable) isUnderDailyTransactions() bool {
	return true
}

func (table FundSuccessTable) isUniqueTransaction() bool {
	return true
}

// func (db FundSuccessTable) getCumulativeAmountFunded(int days) (amtCumulative float64) {
// 	// TODO: get amount from all records for the appropriate days and return it
// 	// e.g. For 1 day, get all fund successes for today(i.e. one second after 23:59:59).
// 	for record := range db {
// 		fmt.Println(record)
// 	}
// }
