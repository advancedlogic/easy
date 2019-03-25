package main

import (
	. "github.com/advancedlogic/easy/easy"
)

func main() {
	microservice, err := NewEasy(
		WithPlugin("plugins/hello.so", "Hello"))

	if err != nil {
		panic(err)
	}

	_, err = microservice.Processor().Process("Upsidedowngalaxy")
	if err != nil {
		microservice.Fatal(err)
	}
}
