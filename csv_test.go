package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsHex(t *testing.T) {

	tests := []struct {
		hex string
		is  bool
	}{
		{"ec", false},
		{"#ececec", true},
		{"ececec", false},
		{"fff", false},
		{"000", false},
		{"#000", true},
		{"#gf2345", false},
	}

	for _, test := range tests {
		assert.True(t, test.is == isHex(test.hex), "Test: %#v", test)
	}
}
