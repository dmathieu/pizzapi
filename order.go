package pizzapi

import (
	"errors"
)

type Order struct {
	Id     int    `json:"id"`
	Item   *Pizza `json:"pizza"`
	Status string `json:"status"`
}

var loadedOrders []*Order

func CreateOrder(item *Pizza) *Order {
	id := len(loadedOrders) + 1
	order := &Order{Id: id, Item: item, Status: "new"}
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
	return loadedOrders, nil
}
