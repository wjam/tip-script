package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleInput(t *testing.T) {
	var tests = []struct {
		input    string
		expected []response
	}{
		{
			"192.168.0.0/23",
			[]response{
				{"text", "192.168.0.0 - 192.168.1.255", ""},
				{"text", "Previous: 192.167.254.0/23", ""},
				{"text", "Next: 192.168.2.0/23", ""},
			},
		},
		{
			"1351700038",
			[]response{{"text", "2012-10-31 16:13:58 +0000 UTC", ""}},
		},
		{
			"1351700038292",
			[]response{{"text", "2012-10-31 16:13:58.292 +0000 UTC", ""}},
		},
		{
			"1351700038292387000",
			[]response{{"text", "2012-10-31 16:13:58.292387 +0000 UTC", ""}},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			actual, err := handleInput(test.input)
			assert.NoError(t, err)
			assert.ElementsMatch(t, test.expected, actual, "actual: %#v", actual)
		})
	}
}
