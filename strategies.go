package main

import (
	"fmt"
)

// ###  strategies
// 1. cautious, adversery votes opposite of initial opinion
// 2. cautious, adversery votes the contrary of the more popular opinion of the last round

// TODO implement interface for the strategies

// - - - - - - - - - - Strategies - - - - - - - - - -
// 1. cautious, adversery votes opposite of initial opinion
func (sim *Sim) getOpinionAdversary_strat1() bool {
	return sim.p.p0 < 0.5
}

// 2. cautious, adversery votes the contrary of the more popular opinion of the last round
func (sim *Sim) getOpinionAdversary_strat2(queriedRound int, advMeasure *AdvMeasure) bool {
	return advMeasure.meanOpinionHonestLastRound < 0.5
}

// - - - - - - - - - - - - - - - - - - - - - - - -

// select which strategy for the adversary
func (sim *Sim) getOpinionAdversary(p *Param, queriedRound, queryingNodeID, queriedNodeID int, advMeasure *AdvMeasure) bool {
	switch p.advStrategy {
	default:
		fmt.Println("Error. Strategy not defined.")
		wait()
	case 1:
		return sim.getOpinionAdversary_strat1()
	case 2:
		return sim.getOpinionAdversary_strat2(queriedRound, advMeasure)
	}
	fmt.Println("Not implemented")
	wait()
	return false
}

// take Measurements to prepare for adversary actions
func (sim *Sim) makeMeasurements(advMeasure *AdvMeasure, thisRound int) {
	switch sim.p.advStrategy {
	default:
		fmt.Println("Not implemented")
		wait()
	case 1:
	case 2:
		sim.measureMeanHonestOpinionLastRound(advMeasure, thisRound-1)
	}
	return
}

func getIndexLastTrue(a []bool) int {
	for i := len(a) - 1; i >= 0; i-- {
		if a[i] == true {
			return i
		}
	}
	fmt.Println("No value True found.")
	wait()
	return -1
}

func getIndexFirstTrue(a []bool) int {
	for i := 0; i < len(a); i++ {
		if a[i] == true {
			return i
		}
	}
	fmt.Println("No value True found.")
	wait()
	return -1
}
