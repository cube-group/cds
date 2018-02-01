package excel

import (
    "testing"
    "fmt"
)

func TestCreate(t *testing.T) {
    data := [][]string{
        {"a1", "b1", "c1"},
        {"d1", "e1", "f1"},
        {"h1", "i1", "j1"},
        {"a2", "b2", "c2"},
        {"d2", "e2", "f2"},
        {"h2", "i2", "j2"},
    }
    err := Create("demoCreate.xlsx", data)
    if err != nil {
        t.Error(fmt.Sprintf("%s", err))
    } else {
        t.Log("ok")
    }
}

func TestRead(t *testing.T) {
    data, err := Read("/Users/chenqionghe/web/trunk/gobase/gobase-trunk/src/alex/excel/demoCreate.xlsx")
    if err != nil {
        t.Error(err)
    }
    fmt.Printf("%#v", data)
}
