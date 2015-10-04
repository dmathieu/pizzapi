package app

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type Pizza struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

var loadedPizzas []*Pizza

func PizzaList() ([]*Pizza, error) {
	if len(loadedPizzas) == 0 {
		file, err := ioutil.ReadFile("pizza.json")
		if err != nil {
			return nil, err
		}
		json.Unmarshal(file, &loadedPizzas)
	}

	return loadedPizzas, nil
}

func FindPizza(id int) (*Pizza, error) {
	pizzas, _ := PizzaList()

	for _, v := range pizzas {
		if v.Id == id {
			return v, nil
		}
	}
	return nil, errors.New("Cannot find pizza " + string(id))
}
