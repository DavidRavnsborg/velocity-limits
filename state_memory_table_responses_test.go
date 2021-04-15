package main

import "testing"

func setupTestResponsesTable() (table ResponsesTable) {
	table = make(ResponsesTable, 0)
	table = append(table, FundResponse{
		Id:         "111",
		CustomerId: "1",
		Accepted:   true,
	})
	table = append(table, FundResponse{
		Id:         "111",
		CustomerId: "2",
		Accepted:   false,
	})
	return table
}

func TestResponsesTableLimitsPositive(t *testing.T) {
	var isUniqueTransaction bool
	table := setupTestResponsesTable()

	isUniqueTransaction = table.isUniqueTransaction("333", "3")
	if !isUniqueTransaction {
		t.Errorf("Expected fund request to be a unique transaction (unique customerId, unique Id).")
	}

	isUniqueTransaction = table.isUniqueTransaction("222", "1")
	if !isUniqueTransaction {
		t.Errorf("Expected fund request to be a unique transaction (non-unique customerId, unique Id).")
	}

	isUniqueTransaction = table.isUniqueTransaction("111", "3")
	if !isUniqueTransaction {
		t.Errorf("Expected fund request to be a unique transaction (unique customerId, non-unique Id).")
	}
}

func TestResponsesTableLimitsNegative(t *testing.T) {
	var isUniqueTransaction bool
	table := setupTestResponsesTable()

	isUniqueTransaction = table.isUniqueTransaction("111", "1")
	if isUniqueTransaction {
		t.Errorf("Expected fund request to not be a unique transaction (non-unique customerId, non-unique Id, accepted=true).")
	}

	isUniqueTransaction = table.isUniqueTransaction("111", "2")
	if isUniqueTransaction {
		t.Errorf("Expected fund request to not be a unique transaction (unique customerId, non-unique Id, accepted=false).")
	}
}
