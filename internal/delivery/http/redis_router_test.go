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
	cases := []struct {
		name   string
		in     *IncrementRequest
		expErr error
	}{
		{
			name: "bad_key",
			in: &IncrementRequest{
				Key: "",
			},
			expErr: ErrEmptyKey,
		},
		{
			name: "bad_value",
			in: &IncrementRequest{
				Key:   "age",
				Value: 0,
			},
			expErr: ErrEmptyIncrementValue,
		},
	}
	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			err := validateIncrementReq(tCase.in)
			require.Error(t, err)
			require.EqualError(t, tCase.expErr, err.Error())
		})
	}
}

func TestValidateSaveRequest(t *testing.T) {
	cases := []struct {
		name   string
		in     *SaveValueRequest
		expErr error
	}{
		{
			name: "valid int value",
			in: &SaveValueRequest{
				Key:   "age",
				Value: 10,
			},
			expErr: nil,
		},
		{
			name: "valid string value",
			in: &SaveValueRequest{
				Key:   "string",
				Value: "string",
			},
			expErr: nil,
		},
	}
	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			err := validateSaveReq(tCase.in)
			require.NoError(t, err)
		})
	}
}

func TestValidateSaveRequestError(t *testing.T) {
	cases := []struct {
		name   string
		in     *SaveValueRequest
		expErr error
	}{
		{
			name: "bad_key",
			in: &SaveValueRequest{
				Key: "",
			},
			expErr: ErrEmptyKey,
		},
	}
	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			err := validateSaveReq(tCase.in)
			require.Error(t, err)
			require.EqualError(t, tCase.expErr, err.Error())
		})
	}
}
