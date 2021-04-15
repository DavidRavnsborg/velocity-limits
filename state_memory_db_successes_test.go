package main

import (
	"testing"
	"time"
)

func setupTestSuccessesDB(customerId string, amounts float64, extraEmptyAmountTransactionsToday int) (db FundSuccessDB) {
	db = make(FundSuccessDB)
	distantPastTimestamp, _ := time.Parse(time.RFC3339, "2000-01-01T00:00:00Z")
	weekStartTimestamp, _ := time.Parse(time.RFC3339, "2021-01-04T00:00:00Z")
	underOneDayTimestamp, _ := time.Parse(time.RFC3339, "2021-01-07T00:00:00Z")
	db[customerId] = append(db[customerId], FundSuccessRecord{
		Id:         "1",
		LoadAmount: amounts,
		Time:       distantPastTimestamp,
	})
	db[customerId] = append(db[customerId], FundSuccessRecord{
		Id:         "2",
		LoadAmount: amounts,
		Time:       weekStartTimestamp,
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

	approved = db["1"].isWithinDailyAmountLimit(fundDepositDate, 1)
	if !approved {
		t.Errorf("Expected fund request to be approved as it is under the daily limit.")
	}

	approved = db["1"].isWithinDailyAmountLimit(fundDepositDate, 4500)
	if !approved {
		t.Errorf("Expected fund request to be approved as it is equal to the daily limit.")
	}

	approved = db["1"].isWithinWeeklyAmountLimit(fundDepositDate, 19000)
	if !approved {
		t.Errorf("Expected fund request to be approved as it is equal to the weekly limit.")
	}

	lastDayOfWeekWithFundDepositsOffset, _ := time.ParseDuration("72h")
	approved = db["1"].isWithinWeeklyAmountLimit(fundDepositDate.Add(lastDayOfWeekWithFundDepositsOffset), 19000)
	if !approved {
		t.Errorf("Expected fund request to be approved as it is equal to the weekly limit.")
	}

	firstDayofWeekWithNoFundDepositsOffset, _ := time.ParseDuration("96h")
	approved = db["1"].isWithinWeeklyAmountLimit(fundDepositDate.Add(firstDayofWeekWithNoFundDepositsOffset), 20000)
	if !approved {
		t.Errorf("Expected fund request to be approved as it is equal to the weekly limit.")
	}

	approved = db["1"].isWithinDailyTransactionLimit(fundDepositDate)
	if !approved {
		t.Errorf("Expected fund request to be approved as it is equal to the daily transaction limit.")
	}
}

func TestFundLimitsDBNegative(t *testing.T) {
	db := setupTestSuccessesDB("1", 500, 2)
	fundDepositDate, _ := time.Parse(time.RFC3339, "2021-01-07T00:13:00Z")
	var approved bool

	approved = db["1"].isWithinDailyAmountLimit(fundDepositDate, 4501)
	if approved {
		t.Errorf("Expected fund request to be declined as it is over the daily limit.")
	}

	approved = db["1"].isWithinWeeklyAmountLimit(fundDepositDate, 19001)
	if approved {
		t.Errorf("Expected fund request to be declined as it is over the weekly limit.")
	}

	lastDayOfWeekWithFundDepositsOffset, _ := time.ParseDuration("72h")
	approved = db["1"].isWithinWeeklyAmountLimit(fundDepositDate.Add(lastDayOfWeekWithFundDepositsOffset), 19001)
	if approved {
		t.Errorf("Expected fund request to be declined as it is over the weekly limit.")
	}

	firstDayofWeekWithNoFundDepositsOffset, _ := time.ParseDuration("96h")
	approved = db["1"].isWithinWeeklyAmountLimit(fundDepositDate.Add(firstDayofWeekWithNoFundDepositsOffset), 20001)
	if approved {
		t.Errorf("Expected fund request to be declined as it is over the weekly limit.")
	}

	approved = db["1"].isWithinDailyTransactionLimit(fundDepositDate)
	if approved {
		t.Errorf("Expected fund request to be declined as it is over the daily transaction limit.")
	}
}
