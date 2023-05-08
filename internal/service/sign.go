package service

import "context"

type SignService struct {
}

func NewSignService() *SignService {
	return &SignService{}
}

func (a *SignService) GenerateSignature(ctx context.Context, text, key string) (string, error) {
	
	return "", nil
}

func (a *SignService) ParseSignature(ctx context.Context, signature, key string) (string, error) {
	return "", nil
}
