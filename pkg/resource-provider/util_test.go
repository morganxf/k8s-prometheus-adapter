package provider

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitDigitLetter(t *testing.T) {
	testCases := []map[string]string{
		{
			"in": "0.0",
			"digit": "0.0",
			"letter": "",
		},
		{
			"in": "0.1",
			"digit": "0.1",
			"letter": "",
		},
		{
			"in": "1",
			"digit": "1",
			"letter": "",
		},
		{
			"in": "0.0m",
			"digit": "0.0",
			"letter": "m",
		},
		{
			"in": "0.1m",
			"digit": "0.1",
			"letter": "m",
		},
		{
			"in": "1m",
			"digit": "1",
			"letter": "m",
		},
		{
			"in": "0.1M",
			"digit": "0.1",
			"letter": "M",
		},
		{
			"in": "1M",
			"digit": "1",
			"letter": "M",
		},
		{
			"in": "1023Ki",
			"digit": "1023",
			"letter": "Ki",
		},
		{
			"in": "1Mi",
			"digit": "1",
			"letter": "Mi",
		},
		{
			"in": "1.1Mi",
			"digit": "1.1",
			"letter": "Mi",
		},
	}
	for _, testCase := range testCases {
		digit, letter := SplitDigitLetter(testCase["in"])
		assert.Equal(t, testCase["digit"], digit)
		assert.Equal(t, testCase["letter"], letter)
	}
}