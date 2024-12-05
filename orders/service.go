package main

import "context"

type Service struct {
	store OrdersStore
}

func NewService(store OrdersStore) *Service {
	return &Service{store}
}

func (s *Service) CreateOrder(ctx context.Context) error {
	return nil
}
