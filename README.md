# simple-types-go

Simple nullable types for Go with first-class support for Postgres and JSON.

simple-types-go provides lightweight wrappers for primitive types (such as `string`, `int`, etc.) that handle NULL values gracefully in both database operations (via database/sql) and JSON marshaling/unmarshaling. These types are useful when dealing with nullable columns in PostgreSQL and when interoperating with APIs where null values are common.

I personally use these types with [sqlc](https://sqlc.dev/) to generate type-safe code for database operations, but they can be used in any Go application that requires nullable types.

## Features

- Supports sql.Scanner and driver.Valuer interfaces for database use (e.g. PostgreSQL).
- Implements json.Marshaler and json.Unmarshaler interfaces.

## Installation

```bash
go get github.com/ItsOnlyGame/simple-types-go
```

## Usage

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/ItsOnlyGame/simple-types-go/types"
)

func main() {
	s := types.NewString("hello")

	data, _ := json.Marshal(s)
	fmt.Println(string(data)) // Output: "hello"

	var s2 types.String
	json.Unmarshal([]byte("null"), &s2)
	fmt.Println(s2.Valid) // Output: false

	fmt.Println(s2.String()) // Output: ""
}
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details
