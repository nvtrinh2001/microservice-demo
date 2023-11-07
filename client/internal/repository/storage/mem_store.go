package storage

import (
	"context"
	"currency/client/internal/entity"
)

type memStore struct {
	db []*entity.Product
}

func NewMemStore(prodList []*entity.Product) *memStore {
	return &memStore{db: prodList}
}

func (ms *memStore) GetProductByID(ctx context.Context, id int) (*entity.Product, error) {
	i := ms.findIndexByProductID(id)
	if id == -1 {
		return nil, entity.ErrProductNotFound
	}

	return ms.db[i], nil
}

func (ms *memStore) findIndexByProductID(id int) int {
	for i, p := range ms.db {
		if p.ID == id {
			return i
		}
	}

	return -1
}
