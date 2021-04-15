package main

import (
	"errors"
	"time"
)

type FundSuccessRecord struct {
	Id         string
	LoadAmount float64
	Time       time.Time
}

type ResponseLimits interface {
	isUniqueTransaction(id string, customerId string) bool
}

type Limits interface {
	updateSuccessAmount(success FundSuccessRecord, customerId string) (err error)
	VelocityLimits
}

type VelocityLimits interface {
	isWithinDailyAmountLimit(requestTime time.Time, amountToAdd float64) bool
	isWithinWeeklyAmountLimit(requestTime time.Time, amountToAdd float64) bool
	isWithinDailyTransactionLimit(requestTime time.Time) bool
}

func checkConditions(userTable Limits, responses ResponseLimits, id string, customerId string, amount float64, time time.Time) (accepted bool, err error) {
	accepted = true

	if !responses.isUniqueTransaction(id, customerId) {
		return false, errors.New(NonUniqueIdError)
	}

	if accepted {
		accepted = userTable.isWithinDailyAmountLimit(time, amount)
	}
	if accepted {
		accepted = userTable.isWithinWeeklyAmountLimit(time, amount)
	}
	if accepted {
		accepted = userTable.isWithinDailyTransactionLimit(time)
	}
	return accepted, nil
}
