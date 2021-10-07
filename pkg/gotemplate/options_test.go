package gotemplate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	min                 = 4
	max                 = 7
	rangeValidatorTests = []struct {
		name        string
		value       int
		expectedErr error
	}{
		{
			name:        "less than min",
			value:       min - 1,
			expectedErr: &ErrOutOfRange{Value: min - 1, Min: min, Max: max},
		},
		{
			name:        "equal to min",
			value:       min,
			expectedErr: nil,
		},
		{
			name:        "equal to max",
			value:       max,
			expectedErr: nil,
		},
		{
			name:        "more than max",
			value:       max + 1,
			expectedErr: &ErrOutOfRange{Value: max + 1, Min: min, Max: max},
		},
	}
)

func Test_RangeValidator(t *testing.T) {
	for _, test := range rangeValidatorTests {
		t.Run(test.name, func(t *testing.T) {
			err := RangeValidator(min, max)(test.value)
			assert.Equal(t, err, test.expectedErr)
		})
	}
}

var (
	regex               = `^this is great$`
	description         = "desc"
	regexValidatorTests = []struct {
		name        string
		value       string
		expectedErr error
	}{
		{
			name:        "no pattern match",
			value:       "this is not so great",
			expectedErr: &ErrInvalidPattern{Value: "this is not so great", Pattern: regex, Description: description},
		},
		{
			name:        "pattern match",
			value:       "this is great",
			expectedErr: nil,
		},
	}
)

func Test_RegexValidator(t *testing.T) {
	for _, test := range regexValidatorTests {
		t.Run(test.name, func(t *testing.T) {
			err := RegexValidator(regex, description)(test.value)
			assert.Equal(t, err, test.expectedErr)
		})
	}
}
