package services

import "github.com/olaoluwavincent/full-course/internal/store"

type PostService struct {
	store store.Storage
}

func NewPostService(store store.Storage) *PostService {
	return &PostService{
		store: store,
	}
}
