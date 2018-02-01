package plugins

import (
    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"
    "alex/log"
    "fmt"
    "runtime"
    "bytes"
)

func MiddleWareError500() martini.Handler {
    return func(c martini.Context, r render.Render, l *log.Logger) {
        defer func() {
            if err := recover(); err != nil {
                fmt.Println("[panic]", err)
                r.HTML(500, "error/error500", string(PanicTrace(5)), render.HTMLOptions{Layout: ""})
            }
        }()

        //c.Next()
    }
}

// PanicTrace trace panic stack info.
func PanicTrace(kb int) []byte {
    s := []byte("/src/runtime/panic.go")
    e := []byte("\ngoroutine ")
    line := []byte("\n")
    stack := make([]byte, kb<<10) //4KB
    length := runtime.Stack(stack, true)
    start := bytes.Index(stack, s)
    stack = stack[start:length]
    start = bytes.Index(stack, line) + 1
    stack = stack[start:]
    end := bytes.LastIndex(stack, line)
    if end != -1 {
        stack = stack[:end]
    }
    end = bytes.Index(stack, e)
    if end != -1 {
        stack = stack[:end]
    }
    stack = bytes.TrimRight(stack, "<br>\n")
    return stack
}
