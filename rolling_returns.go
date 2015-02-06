package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type AnnualRate struct {
	Return float32
	CPI    float32
}

func readReturns() (rates []AnnualRate, err error) { //(years []uint16, returns []float32, cpis []float32, err error) {
	returnsReader, err := os.Open("spx_annual_returns_cpi.txt")
	if err != nil {
		return
	}
	defer returnsReader.Close()
	r := csv.NewReader(returnsReader)
	r.Comma = '\t'
	r.FieldsPerRecord = 3 // Optional, but will return
	records, err := r.ReadAll()
	if err != nil {
		return
	}
	//years = make([]uint16, len(records)-1)
	rates = make([]AnnualRate, len(records)-1)
	for i, v := range records[1:] { // First record is labels so we skip
		//year, _ := strconv.Atoi(v[0]) // Convert string year to int
		//years[i] = uint16(year)
		indexReturn, _ := strconv.ParseFloat(v[1][:len(v[1])-1], 32) // Convert string return to float32, leaves off % at end
		//returns[i] = float32(indexReturn / 100)
		cpi, _ := strconv.ParseFloat(v[2][:len(v[2])-1], 32) // Convert string cpi to float32, leaves off % at end
		//cpis[i] = float32(cpi / 100)
		rates[i] = AnnualRate{Return: float32(indexReturn / 100), CPI: float32(cpi / 100)}
	}
	return
}

func main() {
	//years, returns, cpis, err := readReturns()
	rates, err := readReturns()
	if err != nil {
		fmt.Printf("Error - %s", err)
	}
	// Paramterization
	windows := []int{10, 15, 20, 25}
	consumptions := []float32{.04, .05, .06, .07, .08}
	for _, window := range windows {
		fmt.Printf("For %dy histories:\t", window)
		for _, consumptionBase := range consumptions {
			appreciatedPaCount := 0
			for i := range rates[window:] {
				paBal := float32(1)

				consumption := consumptionBase
				for _, ar := range rates[i : i+window] {
					consumption = consumption * (1 + ar.CPI)
					paBal += paBal*ar.Return - consumption
				}
				if paBal*consumption/consumptionBase > 1 {
					appreciatedPaCount++
				}

			}
			fmt.Printf("%.2f\t", float32(100*appreciatedPaCount)/float32(len(rates[window:])))
		}
		fmt.Println()
	}
}
