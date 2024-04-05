package main

import (
	"testing"

	"github.com/hlubek/readercomp"
)

func TestBatchHandleRequestsToTextFileOutput(t *testing.T) {
	requests := loadData(INPUT_SRC)
	handleBatchRequestsWriteToFile(requests, OUTPUT_TARGET)
	result, err := readercomp.FilesEqual(OUTPUT_TARGET, OUTPUT_SRC)
	if err != nil {
		t.Errorf("Expected equality between the files, received error: %v", err)
	}
	if !result {
		t.Errorf("Expected equality between the files, but FilesEqual returned false.")
	}
}
