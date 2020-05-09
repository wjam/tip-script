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
			"10/8",
			[]response{
				{"text", "10.0.0.0 - 10.255.255.255", ""},
				{"text", "Previous: 9.0.0.0/8", ""},
				{"text", "Next: 11.0.0.0/8", ""},
			},
		},
		{
			"192.168/16",
			[]response{
				{"text", "192.168.0.0 - 192.168.255.255", ""},
				{"text", "Previous: 192.167.0.0/16", ""},
				{"text", "Next: 192.169.0.0/16", ""},
			},
		},
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
