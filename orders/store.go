package main

import "context"

type Store struct {
	// add here our mongoDB
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) Create(ctx context.Context) error {
	return nil
}
