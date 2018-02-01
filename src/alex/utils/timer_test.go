package utils

import (
    "testing"
    "time"
    "fmt"
)

func TestSetInterval(t *testing.T) {
    timer := SetInterval(
        1,
        func(args ...interface{}) {
            fmt.Println(args[0], args[1], args[2])
        },
        "hello",
        "world",
        "golang",
    )
    fmt.Println("setInterval start")
    time.Sleep(time.Second * 10)
    ClearInterval(timer)
    time.Sleep(time.Second * 2)
}

func TestSetTimeout(t *testing.T) {
    timer := SetTimeout(
        1,
        func(args ...interface{}) {
            fmt.Println(args[0], args[1], args[2])
        },
        "hello",
        "world",
        "golang",
    )
    fmt.Println("setTimeout start")
    time.Sleep(time.Second * 10)
    ClearTimeout(timer)
    time.Sleep(time.Second * 2)
}