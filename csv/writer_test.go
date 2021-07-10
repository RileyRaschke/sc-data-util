package csv

import (
	"fmt"
	"testing"
)

func Test_example(t *testing.T) {
	row := CsvRow{}
	if row.String() == "" {
		t.Error("Unable to construct and stringify a CSV row!")
	} else {
		fmt.Print("Stringify'd a CSV Row!\n")
	}
}

/**
* Could be an important one!
 */
func Test_NegativePriceAction(t *testing.T) {
	//row := CsvRow{}
	fmt.Printf("TODO: Write a test for negative prices action! Many birds one stone.\n")
}
