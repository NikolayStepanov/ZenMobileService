package sign

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	log "github.com/sirupsen/logrus"
)

type SignService struct {
}

func NewSignService() *SignService {
	return &SignService{}
}

func (a *SignService) GenerateSignature(ctx context.Context, text, key string) (string, error) {
	hexSignature := ""
	signHash := hmac.New(sha512.New, []byte(key))
	_, err := signHash.Write([]byte(text))
	if err != nil {
		log.Error(err)
		return hexSignature, err
	}
	hexSignature = hex.EncodeToString(signHash.Sum(nil))
	return hexSignature, err
}

func (a *SignService) ValidSignature(ctx context.Context, signature, text, key string) (bool, error) {
	bValidSing := false
	signHash := hmac.New(sha512.New, []byte(key))
	_, err := signHash.Write([]byte(text))
	if err != nil {
		log.Error(err)
		return bValidSing, err
	}
	expectedSign := hex.EncodeToString(signHash.Sum(nil))
	bValidSing = hmac.Equal([]byte(signature), []byte(expectedSign))
	return bValidSing, nil
}
