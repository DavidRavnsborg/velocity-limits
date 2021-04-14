package main

import "time"

type Limits interface {
	VelocityLimits
	OtherLimits
}

type VelocityLimits interface {
	isUnderDailyAmount(requestTime time.Time, amountToAdd float64) bool
	isUnderWeeklyAmount(requestTime time.Time, amountToAdd float64) bool
	isUnderDailyTransactions(requestTime time.Time) bool
}

type OtherLimits interface {
	isUniqueTransaction(id string) bool
}

func checkConditions(userTable Limits, id string, amount float64, time time.Time) (accepted bool) {
	accepted = true
	if accepted {
		accepted = userTable.isUnderDailyAmount(time, amount)
	}
	if accepted {
		accepted = userTable.isUnderWeeklyAmount(time, amount)
	}
	if accepted {
		accepted = userTable.isUnderDailyTransactions(time)
	}
	if accepted {
		accepted = userTable.isUniqueTransaction(id)
	}
	return accepted
}
