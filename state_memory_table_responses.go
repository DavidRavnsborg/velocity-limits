package main

type ResponsesTable []FundResponse

func (responses ResponsesTable) isUniqueTransaction(id string, customerId string) bool {
	for _, response := range responses {
		if response.CustomerId == customerId && response.Id == id {
			return false
		}
	}
	return true
}
