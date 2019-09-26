package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func (p *Param) getCsvParam() []string {
	return []string{
		fmt.Sprintf("%0.9f", p.a),
		fmt.Sprintf("%0.9f", p.deltaab),
		fmt.Sprintf("%0.9f", p.beta),
		fmt.Sprintf("%d", p.k),
		fmt.Sprintf("%d", p.m),
		fmt.Sprintf("%d", p.l),
		fmt.Sprintf("%d", p.n),
		fmt.Sprintf("%0.9f", p.p0),
		fmt.Sprintf("%0.9f", p.q),
		fmt.Sprintf("%d", p.advStrategy),
		fmt.Sprintf("%0.9f", p.rateRandomness),
		fmt.Sprintf("%0.9f", p.deltaWS),
		fmt.Sprintf("%0.9f", p.gammaWS),
		fmt.Sprintf("%d", p.maxTermRound),
	}
}

func initCSV(records [][]string, filename string) error {
	// if the file does not exist create it with the headers
	if !exists("data/result_" + filename + ".csv") {
		f, err := os.Create("data/result_" + filename + ".csv")
		if err != nil {
			fmt.Printf("error creating file: %v", err)
			return err
		}
		defer f.Close()

		w := csv.NewWriter(f)
		w.WriteAll(records) // calls Flush internally

		if err := w.Error(); err != nil {
			log.Fatalln("error writing csv:", err)
		}
		return err
	} else {
		return nil
	}
	panic("Should not be able to get here.")
	return nil // should not get here in any case
}

func writeCSV(records [][]string, filename string, enableHeader bool) error {
	if enableHeader {
		initCSV([][]string{getCsvHeader()}, filename) // requires format into [][]string
	} else {
		initCSV([][]string{}, filename) // requires format into [][]string
	}
	f, err := os.OpenFile("data/result_"+filename+".csv", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("error creating file: %v", err)
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.WriteAll(records) // calls Flush internally

	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}
	return err
}
