package csv

import(
    "testing"
    "fmt"
)

func Test_example( t *testing.T ){
    row := CsvRow{}
    if row.String() == "" {
        t.Error("Unable to construct and stringify a CSV row!")
    } else {
        fmt.Print("Stringify'd a CSV Row!\n")
    }
}

