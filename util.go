package adifparser

import (
	"bytes"
)

// Case-insensitive bytes.Index
func bIndexCI(b, subslice []byte) int {
	return bytes.Index(bytes.ToLower(b), bytes.ToLower(subslice))
}

// Case-insensitive bytes.Contains
func bContainsCI(b, subslice []byte) bool {
	return bytes.Contains(bytes.ToLower(b), bytes.ToLower(subslice))
}

// Find start of next tag
func tagStartPos(b []byte) int {
	nextStart := bytes.IndexByte(b, '<')
	if nextStart == -1 {
		return 0
	}
	return nextStart
}
