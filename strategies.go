package main

import (
	"fmt"
	"sort"
)

// ###  strategies
// 1. cautious, adversery votes opposite of initial opinion
// 2. cautious, adversery votes the contrary of the more popular opinion of the last round
// 3. berserk, adversary tries to delay the process by splitting the opinion
// 4. berserk, adversary tries to delay the process in maximizing the uncertainty
// 5. berserk, try to keep the median of the eta's close to 1/2

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

// 3. berserk: adversary tries to delay the process in splitting the opinions
func (sim *Sim) getOpinionAdversary_strat3(p *Param, queriedRound, queryingNodeID int, advMeasure *AdvMeasure) bool {
	threshold := advMeasure.medianEtaHonestThisRound
	if sim.node[queryingNodeID].etaHonest[queriedRound] > threshold {
		return true
	} else {
		return false
	}
}

// 4. berserk: adversary tries to delay the process in maximizing the uncertainty
func (sim *Sim) getOpinionAdversary_strat4(p *Param, queriedRound, queryingNodeID int, advMeasure *AdvMeasure) bool {
	thresholdLower := p.beta
	thresholdUpper := 1 - p.beta
	if queriedRound == 1 {
		thresholdLower = p.a
		thresholdUpper = p.a + p.deltaab
	}
	threshold := advMeasure.medianEtaHonestThisRound

	if (threshold >= thresholdLower) && (threshold <= thresholdUpper) {
		if sim.node[queryingNodeID].etaHonest[queriedRound] > threshold {
			return true
		} else {
			return false
		}
	} else if threshold <= thresholdLower {
		return true
	} else {
		return false
	}

}

// 5. berserk, try to keep the median of the eta's close to 1/2
func (sim *Sim) prepareOpinionAdversary_strat5(thisRound int, advMeasure *AdvMeasure) {
	etasHelp := make([]float64, sim.p.n_honest)
	attackable := make([]bool, sim.p.n_honest)

	// var etaAttackableHelp []float64
	// var etaAttackableHelpID []float64
	numAttackable := 0
	for id := 0; id < sim.p.n_honest; id++ {
		if !sim.node[id].decided {
			etasHelp[id] = sim.node[id].etaHonest[thisRound]
			attackable[id] = true
			// etaAttackableHelp = append(etaAttackableHelp, sim.node[id].etaHonest[round])
			// etaAttackableHelpID = append(etaAttackableHelpID, id)
			numAttackable++
		} else {
			// also include the opinion of finalized nodes into the median eta
			etasHelp[id] = float64(btoi(sim.node[id].opinion[sim.node[id].TerminationRound]))
			attackable[id] = false
		}
	}
	//create sortable slices
	etas := NewSlice(attackable, etasHelp...)

	for numAttackable > 0 {
		sort.Sort(etas)

		threshold := 0.5
		if thisRound == 1 {
			threshold = sim.p.a + sim.p.deltaab/2
		}
		selectedIndex := 0
		selectedOpinion := false
		// Median below threshold ,
		// TODO make sure the slice is sorted! maybe include a check that the slice is sorted
		if GetUnsortedMedianOfSlice(etas.Float64Slice) < threshold {
			selectedIndex = getIndexLastTrue(etas.flag)
			selectedOpinion = true
			// Median above threshold
		} else {
			selectedIndex = getIndexFirstTrue(etas.flag)
			selectedOpinion = false
		}

		// get node ID
		selectedID := etas.idx[selectedIndex]
		// adv Opinion for this round for the selected node
		advMeasure.CurrentAdvOpinionForNode[selectedID] = selectedOpinion
		pHonest := float64(advMeasure.CurrentNHonestForNode[selectedID]) / float64(sim.p.k)
		etaNew := etas.Float64Slice[selectedIndex]*pHonest + float64(btoi(selectedOpinion))*(1-pHonest)

		// update the eta opinion for this node
		etas.Float64Slice[selectedIndex] = etaNew
		etas.flag[selectedIndex] = false

		// one node less to decide
		numAttackable--
	}
}

func (sim *Sim) getOpinionAdversary_strat5(queryingNodeID int, advMeasure *AdvMeasure) bool {
	return advMeasure.CurrentAdvOpinionForNode[queryingNodeID]
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
	case 3:
		return sim.getOpinionAdversary_strat3(p, queriedRound, queryingNodeID, advMeasure)
	case 4:
		return sim.getOpinionAdversary_strat4(p, queriedRound, queryingNodeID, advMeasure)
	case 5:
		return sim.getOpinionAdversary_strat5(queryingNodeID, advMeasure)
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
	case 3:
		sim.measureMedianEtaHonest(advMeasure, thisRound)
	case 4:
		sim.measureMedianEtaHonest(advMeasure, thisRound)
	case 5:
		sim.prepareOpinionAdversary_strat5(thisRound, advMeasure)
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
