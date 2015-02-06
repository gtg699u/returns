package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func main() {
	returnsReader, err := os.Open("spx_annual_returns_cpi.txt")
	if err != nil {
		return
	}
	defer returnsReader.Close()
	r := csv.NewReader(returnsReader)
	r.Comma = '\t'
	if records, err := r.ReadAll(); err == nil {
		fmt.Printf("First record:\n")
		for i, v := range records[0] {
			fmt.Printf("Field %d:\t%s\n", i, v)
		}
	}
}
