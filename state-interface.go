package main

import (
	"errors"
	"time"
)

type FundSuccessRecord struct {
	Id         string
	LoadAmount float64
	Time       time.Time
	Accepted   bool
}

type ResponseLimits interface {
	isUniqueTransaction(id string, customerId string) bool
}

type Limits interface {
	updateSuccessAmount(success FundSuccessRecord, customerId string) (err error)
	VelocityLimits
}

type VelocityLimits interface {
	isUnderDailyAmount(requestTime time.Time, amountToAdd float64) bool
	isUnderWeeklyAmount(requestTime time.Time, amountToAdd float64) bool
	isUnderDailyTransactions(requestTime time.Time) bool
}

func checkConditions(userTable Limits, responses ResponseLimits, id string, customerId string, amount float64, time time.Time) (accepted bool, err error) {
	accepted = true

	if !responses.isUniqueTransaction(id, customerId) {
		return false, errors.New(NonUniqueIdError)
	}

	if accepted {
		accepted = userTable.isUnderDailyAmount(time, amount)
	}
	if accepted {
		accepted = userTable.isUnderWeeklyAmount(time, amount)
	}
	if accepted {
		accepted = userTable.isUnderDailyTransactions(time)
	}
	return accepted, nil
}
