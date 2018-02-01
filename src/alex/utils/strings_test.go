package utils

import (
    "testing"
    "fmt"
)

func TestStringJoin(t *testing.T) {
    result := StringJoin(" ", "a", "b", "c")
    fmt.Println(result)
}
