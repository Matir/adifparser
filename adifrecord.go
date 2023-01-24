package adifparser

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// Public interface for ADIFRecords
type ADIFRecord interface {
	// Print as ADIF String
	ToString() string
	// Fingerprint for duplication detection
	Fingerprint() string
	// Setters and getters
	GetValue(string) (string, error)
	SetValue(string, string)
	// Get all of the present field names
	GetFields() []string
}

// Internal implementation for ADIFRecord
type baseADIFRecord struct {
	values map[string]string
}

type fieldData struct {
	name     string
	value    string
	typecode byte
	hasType  bool
}

// Errors
var NoData = errors.New("No data to parse.")
var NoSuchField = errors.New("No such field.")
var InvalidField = errors.New("Invalid field definition.")

// Create a new ADIFRecord from scratch
func NewADIFRecord() *baseADIFRecord {
	record := &baseADIFRecord{}
	record.values = make(map[string]string)
	return record
}

func serializeField(name string, value string) string {
	return fmt.Sprintf("<%s:%d>%s", name, len(value), value)
}

// Print an ADIFRecord as a string
func (r *baseADIFRecord) ToString() string {
	var record bytes.Buffer
	for _, n := range ADIFfieldOrder {
		if v, ok := r.values[n]; ok {
			record.WriteString(serializeField(n, v))
		}
	}
	// Handle custom fields
	for n, v := range r.values {
		if !isStandardADIFField(n) {
			record.WriteString(serializeField(n, v))
		}
	}
	return record.String()
}

// Get fingerprint of ADIFRecord
func (r *baseADIFRecord) Fingerprint() string {
	fpfields := []string{
		"call", "station_callsign", "band",
		"freq", "mode", "qso_date", "time_on",
		"time_off"}
	fpvals := make([]string, 0, len(fpfields))
	for _, f := range fpfields {
		if n, ok := r.values[f]; ok {
			fpvals = append(fpvals, n)
		}
	}
	fptext := strings.Join(fpvals, "|")
	h := sha256.New()
	h.Write([]byte(fptext))
	return hex.EncodeToString(h.Sum(nil))
}

// Get a value
func (r *baseADIFRecord) GetValue(name string) (string, error) {
	if v, ok := r.values[name]; ok {
		return v, nil
	}
	return "", NoSuchField
}

// Set a value
func (r *baseADIFRecord) SetValue(name string, value string) {
	r.values[strings.ToLower(name)] = value
}

// Get all of the present field names
func (r *baseADIFRecord) GetFields() []string {
	keys := make([]string, len(r.values))
	i := 0
	for k := range r.values {
		keys[i] = k
		i++
	}
	return keys
}

// Delete a field (from the internal map)
func (r *baseADIFRecord) DeleteField(name string) (bool, error) {
	if _, ok := r.values[name]; ok {
		delete(r.values, name)
		return true, nil
	}
	return false, NoSuchField
}
