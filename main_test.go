package main

import (
	"testing"

	"github.com/hlubek/readercomp"
)

func TestBatchHandleRequestsToTextFileOutput(t *testing.T) {
	requests := loadData()
	handleBatchRequests(requests)
	result, err := readercomp.FilesEqual(output_target, output_src)
	if err != nil {
		t.Errorf("Expected equality between the files, received error: %v", err)
	}
	if !result {
		t.Errorf("Expected equality between the files, but FilesEqual returned false.")
	}
}
