package util

import (
    "testing"
    "fmt"
    "os"
)

func Test_utilTests(t *testing.T){
    b, err := WriteBuffer(os.Stdout)
    if err != nil {
        t.Errorf("Failed to create a write buffer from os.Stdout!")
    }
    b.Write([]byte("*os.File -> buffer!\n"))
    fmt.Printf("WriteBuffer(*os.File) works!\n")
}
