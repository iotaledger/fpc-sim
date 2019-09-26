package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func parseInput(filename string) []*Param {
	//initiate param slice
	P := make([]*Param, 0)
	p := &Param{
		EnableQueryWithRepetition: true,
	}
	P = append(P, p)

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	s := bufio.NewScanner(f)
	for s.Scan() {
		if !strings.Contains(s.Text(), "#") && len(s.Text()) > 0 {
			param := strings.Split(s.Text(), "=")
			val := strings.Split(param[1], ",")
			//store the current set of parameter sets
			P0 := make([]*Param, len(P))
			copy(P0, P)
			for iValue, value := range val {
				//add another set of P0 if it's not the first parameter
				if iValue > 0 {
					for _, p0 := range P0 {
						p1 := *p0 //copy values without pointer reference
						P = append(P, &p1)
					}
				}
				for iP0 := range P0 { // for each of the elements in P0 update this new value
					P[iValue*len(P0)+iP0].selectParam(param[0], value, err)
				}
			}
		}
	}
	fmt.Println()
	err = s.Err()
	if err != nil {
		log.Fatal(err)
	}

	return P
}

//evaluate text file
func (p *Param) selectParam(paramtype, value string, err error) {
	switch paramtype {
	case "a":
		p.a, err = strconv.ParseFloat(value, 64)
		if err != nil {
			panic("Unable to parse a")
		}
	case "deltaab":
		p.deltaab, err = strconv.ParseFloat(value, 64)
		if err != nil {
			panic("Unable to parse deltaab")
		}
	case "beta":
		p.beta, err = strconv.ParseFloat(value, 64)
		if err != nil {
			panic("Unable to parse beta")
		}
	case "p0":
		p.p0, err = strconv.ParseFloat(value, 64)
		if err != nil {
			panic("Unable to parse p0")
		}
	case "q":
		p.q, err = strconv.ParseFloat(value, 64)
		if err != nil {
			panic("Unable to parse q")
		}
	case "rateRandomness":
		p.rateRandomness, err = strconv.ParseFloat(value, 64)
		if err != nil {
			panic("Unable to parse rateRandomness")
		}
	case "nRun":
		Nrun, err := strconv.ParseInt(value, 10, 64)
		p.Nrun = int(Nrun)
		if err != nil {
			panic("Unable to parse nRun")
		}
	case "N":
		N, err := strconv.ParseInt(value, 10, 64)
		p.n = int(N)
		if err != nil {
			panic("Unable to parse N")
		}
	case "k":
		k, err := strconv.ParseInt(value, 10, 64)
		p.k = int(k)
		if err != nil {
			panic("Unable to parse k")
		}
	case "m":
		m, err := strconv.ParseInt(value, 10, 64)
		p.m = int(m)
		if err != nil {
			panic("Unable to parse m")
		}
	case "l":
		l, err := strconv.ParseInt(value, 10, 64)
		p.l = int(l)
		if err != nil {
			panic("Unable to parse l")
		}
	case "maxTermRound":
		maxTermRound, err := strconv.ParseInt(value, 10, 64)
		p.maxTermRound = int(maxTermRound)
		if err != nil {
			panic("Unable to parse maxTermRound")
		}
	case "strategy":
		advStrategy, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			panic("Unable to parse strategy")
		}
		p.advStrategy = int(advStrategy)
	case "enableSaveEta":
		if value == "false" {
			p.enableSaveEta = false
		} else if value == "true" {
			p.enableSaveEta = true
		} else {
			panic("Unable to parse enableSaveEta")
		}
	case "enableRandN_adv":
		if value == "false" {
			p.enableRandN_adv = false
		} else if value == "true" {
			p.enableRandN_adv = true
		} else {
			panic("Unable to parse enableRandN_adv")
		}
	case "enableExtremeBeta":
		if value == "false" {
			p.enableExtremeBeta = false
		} else if value == "true" {
			p.enableExtremeBeta = true
		} else {
			panic("Unable to parse enableExtremeBeta")
		}
	case "etaAgreement":
		p.etaAgreement, err = strconv.ParseFloat(value, 64)
		if err != nil {
			panic("Unable to parse etaAgreement")
		}
	case "enableWS":
		if value == "false" {
			p.enableWS = false
		} else if value == "true" {
			p.enableWS = true
		} else {
			panic("Unable to parse enableWS")
		}
	case "gammaWS":
		p.gammaWS, err = strconv.ParseFloat(value, 64)
		if err != nil {
			panic("Unable to parse gammaWS")
		}
	case "deltaWS":
		p.deltaWS, err = strconv.ParseFloat(value, 64)
		if err != nil {
			panic("Unable to parse deltaWS")
		}
	}

}
