package storage

import "errors"

type StorageFake struct {
	data map[string]string
}

func NewFake() *StorageFake {
	return &StorageFake{
		data: make(map[string]string),
	}
}

func (s *StorageFake) Create(orderUid string, data string) error {
	if _, ok := s.data[orderUid]; ok {
		return errors.New("already exist")
	}
	s.data[orderUid] = data
	return nil
}

func (s *StorageFake) Get(uid string) (string, error) {
	if _, ok := s.data[uid]; !ok {
		return "", errors.New("not found")
	}
	return s.data[uid], nil
}

func (s *StorageFake) GetAll() []OrderPair {
	result := make([]OrderPair, 0)
	for k, v := range s.data {
		result = append(result, OrderPair{
			OrderUID: k,
			Data:     v,
		})
	}
	return result
}
