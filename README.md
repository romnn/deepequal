## deepequal

[![GitHub](https://img.shields.io/github/license/romnn/deepequal)](https://github.com/romnn/deepequal)
[![GoDoc](https://godoc.org/github.com/romnn/deepequal?status.svg)](https://godoc.org/github.com/romnn/deepequal)
[![Test Coverage](https://codecov.io/gh/romnn/deepequal/branch/master/graph/badge.svg)](https://codecov.io/gh/romnn/deepequal)
[![Release](https://img.shields.io/github/release/romnn/deepequal)](https://github.com/romnn/deepequal/releases/latest)

This package is based on the original `reflect.DeepEqual`, but adds useful error messages pointing out where and how the compared values differ.

```go
import "github.com/romnn/deepequal"
```

#### Example
```go
// examples/example1/main.go

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
	if equal, err := deepequal.DeepEqual(a, b); !equal {
		log.Fatalf("not equal: %v", err)
	}
}

```

For more examples see the `examples/` directory.

#### Acknowledgement
- Check out the vanilla `reflect` implementation at [golang.org/src/reflect/deepequal.go](https://golang.org/src/reflect/deepequal.go)
