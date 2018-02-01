package cli

import (
    "testing"
    "fmt"
)

func TestCli_Run(t *testing.T) {
    cli := New("123.57.157.78", "fyadmin", "fy@admin$1&#", 10068)
    output, err := cli.Run("free -h")
    fmt.Printf("%v\n%v", output, err)

}
