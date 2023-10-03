package order

import (
	"errors"
	"l0/internal/service/storage"
)

type Storage interface {
	Create(orderUid string, data string) error
	Get(uid string) (string, error)
	GetAll() []storage.OrderPair
}

type CacheStorage interface {
	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
}

type Service struct {
	storage Storage
	cache   CacheStorage
}

func New(storage Storage, cache CacheStorage) *Service {

	orders := storage.GetAll()
	for _, order := range orders {
		cache.Set(order.OrderUID, order.Data)
	}

	return &Service{
		storage: storage,
		cache:   cache,
	}
}

func (svc *Service) Create(orderUID, data string) {
	if err := svc.storage.Create(orderUID, data); err != nil {
		return
	}
	svc.cache.Set(orderUID, data)
}

func (svc *Service) Get(uid string) (string, error) {
	cacheData, ok := svc.cache.Get(uid)
	if ok {
		return cacheData.(string), nil
	}

	data, err := svc.storage.Get(uid)
	if err == nil {
		return data, nil
	}

	return "", errors.New("not found")
}
