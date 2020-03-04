## deepequal

This package is based on the original `reflect.DeepEqual` check, but adds useful error messages pointing out where and how the compared values violate the equality constraint.

```go
import "github.com/romnnn/deepequal"
```

#### Example
```go
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
```

For more examples see the `examples/` directory. To see an example in action, you can also run it:
```bash
go run github.com/romnnn/deepequal/examples/example1
```

#### Credits
- Check out the vanilla `reflect` implementation at [golang.org/src/reflect/deepequal.go](https://golang.org/src/reflect/deepequal.go)