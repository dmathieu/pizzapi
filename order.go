package pizzapi

import (
	"errors"
	"time"
)

type orderType struct {
	Name     string
	SetAfter time.Duration
}

var orderTypes = []*orderType{
	&orderType{"new", 0},
	&orderType{"processing", 1 * time.Minute},
	&orderType{"cooking", 3 * time.Minute},
	&orderType{"delivering", 10 * time.Minute},
	&orderType{"finished", 15 * time.Minute},
}

type Order struct {
	Id        int    `json:"id"`
	Item      *Pizza `json:"pizza"`
	Status    string `json:"status"`
	CreatedAt time.Time
}

var loadedOrders []*Order

func CreateOrder(item *Pizza) *Order {
	id := len(loadedOrders) + 1
	order := &Order{
		Id:        id,
		Item:      item,
		Status:    orderTypes[0].Name,
		CreatedAt: time.Now(),
	}
	loadedOrders = append(loadedOrders, order)

	return order
}

func FindOrder(id int) (*Order, error) {
	for _, v := range loadedOrders {
		if v.Id == id {
			return v, nil
		}
	}
	return nil, errors.New("Cannot find order " + string(id))
}

func AllOrders() ([]*Order, error) {
	if loadedOrders == nil {
		return []*Order{}, nil
	}
	return loadedOrders, nil
}
