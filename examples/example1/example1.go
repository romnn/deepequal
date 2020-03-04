package main

import (
	"log"

	"github.com/romnnn/deepequal"
)

type person struct {
	Name    string
	Age     int
	Hobbies []string
}

func main() {
	a := person{Name: "A", Age: 22, Hobbies: []string{"Surfing"}}
	b := person{Name: "A", Age: 22, Hobbies: []string{}}
	if equal, err := deepequal.DeepEqual(a, b); !equal {
		log.Fatalf("Not DeepEqual because of: %s", err.Error())
	}
}
