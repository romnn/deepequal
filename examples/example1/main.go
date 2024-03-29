package main

import (
	"log"

	"github.com/romnn/deepequal"
)

type person struct {
	Name    string
	Age     int
	Hobbies []string
}

func main() {
	a := person{Name: "A", Age: 22, Hobbies: []string{"Surfing"}}
	b := person{Name: "A", Age: 22, Hobbies: []string{}}
	if equal, err := deepequal.DeepEqual(a, a); !equal {
		log.Fatalf("not equal: %v", err)
	}

	if equal, err := deepequal.DeepEqual(a, b); equal {
		log.Fatalf("unexpected equal: %v", err)
	}
}
