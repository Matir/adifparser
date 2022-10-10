package adifparser

import (
	"bytes"
)

// ASCII lowercase converter
// For a byte as an ASCII character
func charToLower(c byte) byte {
	if 'A' <= c && c <= 'Z' {
		c += 'a' - 'A'
	}
	return c
}

// Strictly-ASCII-only lowercase converter
// For a byte sequence
// No Unicode processing
// See bytes.ToLower() source code
func bStrictToLower(s []byte) []byte {
	hasUpper := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		hasUpper = hasUpper || ('A' <= c && c <= 'Z')
	}
	if !hasUpper {
		return append([]byte(""), s...)
	}
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if 'A' <= c && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return b
}

// ASCII uppercase converter
// For a byte as an ASCII character
func charToUpper(c byte) byte {
	if 'a' <= c && c <= 'z' {
		c -= 'a' - 'A'
	}
	return c
}

// Strictly-ASCII-only uppercase converter
// For a byte sequence
// No Unicode processing
// See bytes.ToUpper() source code
func bStrictToUpper(s []byte) []byte {
	hasLower := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		hasLower = hasLower || ('a' <= c && c <= 'z')
	}
	if !hasLower {
		return append([]byte(""), s...)
	}
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if 'a' <= c && c <= 'z' {
			c -= 'a' - 'A'
		}
		b[i] = c
	}
	return b
}

// Case-insensitive bytes.Index
// This function only handles ASCII bytes - no Unicode-specific conversion
func bIndexCI(b, subslice []byte) int {
	return bytes.Index(bStrictToLower(b), bStrictToLower(subslice))
}

// Case-insensitive bytes.Contains
// This function only handles ASCII bytes - no Unicode-specific conversion
func bContainsCI(b, subslice []byte) bool {
	return bytes.Contains(bStrictToLower(b), bStrictToLower(subslice))
}

// Find start of next tag
func tagStartPos(b []byte) int {
	nextStart := bytes.IndexByte(b, '<')
	if nextStart == -1 {
		return 0
	}
	return nextStart
}
