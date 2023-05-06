package http

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateIncrementRequest(t *testing.T) {
	reqIncr := &IncrementRequest{
		Key:   "age",
		Value: 10,
	}
	err := validateIncrementReq(reqIncr)
	require.NoError(t, err)
}

func TestValidateIncrementRequestError(t *testing.T) {
	tests := []struct {
		name          string
		input         *IncrementRequest
		expectedError error
	}{
		{
			name: "bad_key",
			input: &IncrementRequest{
				Key: "",
			},
			expectedError: ErrEmptyKey,
		},
		{
			name: "bad_value",
			input: &IncrementRequest{
				Key:   "age",
				Value: 0,
			},
			expectedError: ErrEmptyIncrementValue,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateIncrementReq(test.input)
			require.Error(t, err)
			require.EqualError(t, test.expectedError, err.Error())
		})
	}
}

func TestValidateSaveRequest(t *testing.T) {
	tests := []struct {
		name          string
		input         *SaveValueRequest
		expectedError error
	}{
		{
			name: "valid int value",
			input: &SaveValueRequest{
				Key:   "age",
				Value: 10,
			},
			expectedError: nil,
		},
		{
			name: "valid string value",
			input: &SaveValueRequest{
				Key:   "string",
				Value: "string",
			},
			expectedError: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateSaveReq(test.input)
			require.NoError(t, err)
		})
	}
}

func TestValidateSaveRequestError(t *testing.T) {
	tests := []struct {
		name          string
		input         *SaveValueRequest
		expectedError error
	}{
		{
			name: "bad_key",
			input: &SaveValueRequest{
				Key: "",
			},
			expectedError: ErrEmptyKey,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateSaveReq(test.input)
			require.Error(t, err)
			require.EqualError(t, test.expectedError, err.Error())
		})
	}
}
