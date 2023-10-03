package stan

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/nats-io/stan.go"
)

func orderIsValid(order *Order) bool {
	validate := validator.New()
	err := validate.Struct(order)
	return err == nil
}

func (s *Service) handleMessage(m *stan.Msg) {
	var order Order
	if err := json.Unmarshal(m.Data, &order); err != nil {
		return
	}

	if !orderIsValid(&order) {
		return
	}

	dataBytes, err := json.Marshal(order)
	if err != nil {
		return
	}
	data := string(dataBytes)

	s.orderSvc.Create(order.OrderUID, data)
}
