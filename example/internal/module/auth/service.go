package auth

import (
	"context"
	"encoding/binary"
	"errors"

	"github.com/go-chocolate/chocolate/pkg/kv"
	"github.com/google/uuid"
)

type Service interface {
	Token(ctx context.Context, id int64) (string, error)
	Validate(ctx context.Context, token string) (int64, error)
}

type UnimplementedService struct{}

func (s *UnimplementedService) Token(ctx context.Context, id int64) (string, error) {
	return "", errors.New("unimplemented")
}
func (s *UnimplementedService) Validate(ctx context.Context, token string) (int64, error) {
	return 0, errors.New("unimplemented")
}

func NewService(kv kv.Storage) Service {
	return &service{kv: kv}
}

type service struct {
	kv kv.Storage
}

func (s *service) Token(ctx context.Context, id int64) (string, error) {
	token := uuid.New().String()
	var b = make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(id))
	if err := s.kv.Set(ctx, token, b); err != nil {
		return "", err
	}
	return token, nil
}

func (s *service) Validate(ctx context.Context, token string) (int64, error) {
	val, err := s.kv.Get(ctx, token)
	if err != nil {
		return 0, err
	}
	if len(val) != 8 {
		return 0, errors.New("invalid token")
	}
	return int64(binary.BigEndian.Uint64(val)), err

}
