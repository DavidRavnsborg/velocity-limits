package main

import (
	"testing"
	"time"
)

// func setupTestDB(customerId string, distantPastTimestamp time.Time, underSevenDaysTimestamp time.Time, underOneDayTimestamp time.Time) (db FundSuccessDB) {
func setupTestSuccessesDB(customerId string, amounts float64, extraEmptyAmountTransactionsToday int) (db FundSuccessDB) {
	db = make(FundSuccessDB)
	distantPastTimestamp, _ := time.Parse(time.RFC3339, "2000-01-01T00:00:00Z")
	underSevenDaysTimestamp, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	underOneDayTimestamp, _ := time.Parse(time.RFC3339, "2021-01-07T00:00:00Z")
	db[customerId] = append(db[customerId], FundSuccessRecord{
		Id:         "1",
		LoadAmount: amounts,
		Time:       distantPastTimestamp,
	})
	db[customerId] = append(db[customerId], FundSuccessRecord{
		Id:         "2",
		LoadAmount: amounts,
		Time:       underSevenDaysTimestamp,
	})
	db[customerId] = append(db[customerId], FundSuccessRecord{
		Id:         "3",
		LoadAmount: amounts,
		Time:       underOneDayTimestamp,
	})
	for i := 0; i < extraEmptyAmountTransactionsToday; i++ {
		db[customerId] = append(db[customerId], FundSuccessRecord{
			Id:         "X",
			LoadAmount: 0,
			Time:       underOneDayTimestamp,
		})
	}
	return db
}

func TestFundLimitsDBPositive(t *testing.T) {
	db := setupTestSuccessesDB("1", 500, 1)
	fundDepositDate, _ := time.Parse(time.RFC3339, "2021-01-07T00:13:00Z")
	var approved bool

	approved = db["1"].isUnderDailyAmount(fundDepositDate, 1)
	if !approved {
		t.Errorf("Expected fund request to be approved as it is under to the daily limit.")
	}

	approved = db["1"].isUnderDailyAmount(fundDepositDate, 4500)
	if !approved {
		t.Errorf("Expected fund request to be approved as it is equal to the daily limit.")
	}

	approved = db["1"].isUnderWeeklyAmount(fundDepositDate, 19000)
	if !approved {
		t.Errorf("Expected fund request to be approved as it is equal to the weekly limit.")
	}

	approved = db["1"].isUnderDailyTransactions(fundDepositDate)
	if !approved {
		t.Errorf("Expected fund request to be approved as it is equal to the daily transaction limit.")
	}
}

func TestFundLimitsDBNegative(t *testing.T) {
	db := setupTestSuccessesDB("1", 500, 2)
	fundDepositDate, _ := time.Parse(time.RFC3339, "2021-01-07T00:13:00Z")
	var approved bool

	approved = db["1"].isUnderDailyAmount(fundDepositDate, 4501)
	if approved {
		t.Errorf("Expected fund request to be declined as it is over the daily limit.")
	}

	approved = db["1"].isUnderWeeklyAmount(fundDepositDate, 19001)
	if approved {
		t.Errorf("Expected fund request to be declined as it is over the weekly limit.")
	}

	approved = db["1"].isUnderDailyTransactions(fundDepositDate)
	if approved {
		t.Errorf("Expected fund request to be declined as it is over the daily transaction limit.")
	}
}
