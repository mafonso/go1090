package main

import (
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {

	message := []byte{0x24, 0x4b, 0xbb, 0x9a, 0xc9, 0xf0}
	utcDate := time.Now().UTC()

	wanted := time.Date(utcDate.Year(), utcDate.Month(), utcDate.Day(),
		10, 19, 26, 999999984, time.UTC)

	result := parseTime(message)

	if result != wanted {
		t.Errorf("Time not parsed correctly, got: %s, want: %s.", result, wanted)
	}
}
