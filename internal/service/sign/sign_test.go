package sign

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateGenerateAndValidSignature(t *testing.T) {
	tests := []struct {
		name          string
		inputText     string
		inputKey      string
		outputSign    string
		validSign     bool
		expectedError error
	}{
		{
			name:       "Test1Success",
			inputText:  "test",
			inputKey:   "test123",
			validSign:  true,
			outputSign: "b596e24739fd44d42ffd25f26ea367dad3a71f61c8c5fab6b6ee6ceeae5a7170b66445d6eaadfb49e6d4e968a2888726ff522e3bf065c966aa66a24153778382",
		},
		{
			name:       "Test2Success",
			inputText:  "test",
			inputKey:   "test",
			validSign:  true,
			outputSign: "9ba1f63365a6caf66e46348f43cdef956015bea997adeb06e69007ee3ff517df10fc5eb860da3d43b82c2a040c931119d2dfc6d08e253742293a868cc2d82015",
		},
		{
			name:       "Test3Success",
			inputText:  "TestTest123",
			inputKey:   "123",
			validSign:  true,
			outputSign: "a83745f7e38b471856389faad8c4def134181c8a83fb592e158b82ac3d39dc2990fa021dfa2b45c31e76fc17e45131049a900945e2ae53ec4bc9c2069f9f10b9",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			singService := NewSignService()
			ctx := context.Background()
			sign, err := singService.GenerateSignature(ctx, test.inputText, test.inputKey)
			validSign, err := singService.ValidSignature(ctx, sign, test.inputText, test.inputKey)
			assert.Equal(t, test.validSign, validSign)
			require.NoError(t, err)
			assert.Equal(t, test.outputSign, sign)
		})
	}
}

func TestValidateValidSignature(t *testing.T) {
	tests := []struct {
		name          string
		inputText     string
		inputKey      string
		outputSign    string
		validSign     bool
		expectedError error
	}{
		{
			name:       "Test1Success",
			inputText:  "test",
			inputKey:   "test123",
			validSign:  true,
			outputSign: "b596e24739fd44d42ffd25f26ea367dad3a71f61c8c5fab6b6ee6ceeae5a7170b66445d6eaadfb49e6d4e968a2888726ff522e3bf065c966aa66a24153778382",
		},
		{
			name:       "Test2WrongText",
			inputText:  "test2",
			inputKey:   "test",
			validSign:  false,
			outputSign: "9ba1f63365a6caf66e46348f43cdef956015bea997adeb06e69007ee3ff517df10fc5eb860da3d43b82c2a040c931119d2dfc6d08e253742293a868cc2d82015",
		},
		{
			name:       "Test3WrongKey",
			inputText:  "TestTest123",
			inputKey:   "123h",
			validSign:  false,
			outputSign: "a83745f7e38b471856389faad8c4def134181c8a83fb592e158b82ac3d39dc2990fa021dfa2b45c31e76fc17e45131049a900945e2ae53ec4bc9c2069f9f10b9",
		},
		{
			name:       "Test4WrongOutputSign",
			inputText:  "TestTest123",
			inputKey:   "123h",
			validSign:  false,
			outputSign: "a83745f7e38b471856389faad8c4de",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			singService := NewSignService()
			ctx := context.Background()
			validSign, err := singService.ValidSignature(ctx, test.outputSign, test.inputText, test.inputKey)
			assert.Equal(t, test.validSign, validSign)
			require.NoError(t, err)
		})
	}
}
