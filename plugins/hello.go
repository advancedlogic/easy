package main

import (
	"fmt"
	"github.com/advancedlogic/easy/interfaces"
)

type hello struct{}

func (h hello) Init(service interfaces.Service) error {
	return nil
}

func (h hello) Close() error {
	return nil
}

func (h hello) Process(data interface{}) (interface{}, error) {
	fmt.Println(data.(string))
	return data, nil
}

var Hello hello
