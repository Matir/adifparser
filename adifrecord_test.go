package adifparser

import (
	"strings"
	"testing"
)

func TestGetFields(t *testing.T) {
	testData := "<call:4>W1AW<STATION_CALL:6>KF4MDV <EOR>"
	expected := [2]string{"call", "station_call"}
	buf := strings.NewReader(testData)
	reader := NewADIFReader(buf)

	record, err := reader.ReadRecord()
	if err != nil {
		t.Fatal(err)
	}
	fieldNames := record.GetFields()

	if len(fieldNames) != len(expected) {
		t.Fatalf("Expected %d fields but got %d", len(expected), len(fieldNames))
	}

OUTER:
	for _, exp := range expected {
		for _, field := range fieldNames {
			if exp == field {
				continue OUTER
			}
		}
		t.Fatalf("Expected field %v wasn't in the actual fields", exp)
	}
}

func TestSetValue(t *testing.T) {
	testData := map[string]string{"call": "W1AW", "STATION_CALL": "KF4MDV"}
	expected := [2]string{"call", "station_call"}

	record := NewADIFRecord()
	for k, v := range testData {
		record.SetValue(k, v)
	}
	fieldNames := record.GetFields()

	if len(fieldNames) != len(expected) {
		t.Fatalf("Expected %d fields but got %d", len(expected), len(fieldNames))
	}

OUTER:
	for _, exp := range expected {
		for _, field := range fieldNames {
			if exp == field {
				continue OUTER
			}
		}
		t.Fatalf("Expected field %v wasn't in the actual fields", exp)
	}
}
