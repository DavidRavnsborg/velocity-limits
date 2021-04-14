package main

import "time"

type Limits interface {
	VelocityLimits
	OtherLimits
}

type VelocityLimits interface {
	isUnderDailyAmount() bool
	isUnderWeeklyAmount() bool
	isUnderDailyTransactions() bool
}

type OtherLimits interface {
	isUniqueTransaction() bool
}

func checkConditions(userTable Limits, id string, amount float64, time time.Time) (accepted bool) {
	accepted = true
	if accepted {
		accepted = userTable.isUnderDailyAmount()
	}
	if accepted {
		accepted = userTable.isUnderWeeklyAmount()
	}
	if accepted {
		accepted = userTable.isUnderDailyTransactions()
	}
	if accepted {
		accepted = userTable.isUnderWeeklyAmount()
	}
	if accepted {
		accepted = userTable.isUniqueTransaction()
	}
	return accepted
}
