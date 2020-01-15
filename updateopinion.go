package main

import "fmt"

// get random threshold for this round
func (sim *Sim) getU(threshold_l, threshold_u float64) float64 {
	if sim.rand.Float64() <= sim.p.rateRandomness {
		// discrete random : upper and lower threshold
		if sim.p.enableExtremeBeta {
			if sim.rand.Intn(2) == 0 {
				return threshold_u
			} else {
				return threshold_l
			}
			// continuous random
		} else {
			return threshold_l + sim.rand.Float64()*(threshold_u-threshold_l)
		}
	} else {
		return (threshold_u + threshold_l) / 2
	}
}

// query p.k random nodes, self query and repetition are allowed
func (sim *Sim) querySampleHonest(sample []int, queriedRound, queryingNodeID, k int) (float64, float64) {
	eta := 0.
	counterHonest := 0.
	for _, queriedNodeID := range sample {
		if queriedNodeID < sim.p.n_honest {
			counterHonest++
			actualRound := queriedRound
			if sim.node[queriedNodeID].decided {
				actualRound = sim.node[queriedNodeID].TerminationRound
			}
			eta += float64(btoi(sim.node[queriedNodeID].opinion[actualRound]))
		}
	}
	// return etaHonest, pHonest
	if counterHonest > 0 {
		return eta / counterHonest, counterHonest / float64(k)
	} else {
		return 0., 0.
	}
}

func (sim *Sim) querySampleAdversary(p *Param, sample []int, queriedRound, queryingNodeID int, advMeasure *AdvMeasure) float64 {
	// provide opinions
	eta := 0.
	counterAdversary := 0.
	for _, queriedNodeID := range sample {
		if queriedNodeID >= p.n_honest {
			counterAdversary++
			eta += float64(btoi(sim.getOpinionAdversary(p, queriedRound, queryingNodeID, queriedNodeID, advMeasure)))
		}
	}
	// return etaAdversary
	if counterAdversary > 0 {
		return eta / counterAdversary
	} else {
		return 0.
	}
}

// get sample
func (sim *Sim) getSample(id, k int) []int {
	sample := make([]int, k)
	if sim.p.EnableQueryWithRepetition {
		if !sim.p.enableZipf {
			// no mana distribution
			for i := 0; i < k; i++ {
				sample[i] = sim.node[id].neighborsID[sim.rand.Intn(len(sim.node[id].neighborsID))]
			}
		} else {
			// with mana distribution
			for i := 0; i < k; i++ {
				randfloat := sim.rand.Float64()
				id2 := 0
				for ; sim.node[id].neighborsCumuManaVec[id2] < randfloat; id2++ {
				}
				sample[i] = sim.node[id].neighborsID[id2]
			}
		}
		// sim.node[i].opinion = append(sim.node[i].opinion, sim.querysample(p, round-1, i) > U)
	} else {
		panic("This has to be checked before enabling it - e.g. it does not include the neighborsIDlist")
		selected := make(map[int]bool)
		// sample k unique random nodes
		if k > sim.p.n/2 {
			fmt.Println("It is not recommended to use QueryWithNoRepetition for k>n/2")
			wait()
		}
		for len(selected) < k {
			selected[sim.rand.Intn(sim.p.n)] = true // create an entry for that ID
		}
		fmt.Println(selected)
		counter := 0
		for ID, _ := range selected {
			counter++
			sample[counter-1] = ID
		}
		fmt.Println("Check that this last loop works")
		wait()
		// sim.node[i].opinion = append(sim.node[i].opinion, sim.querySampleUnique(p, round-1, i) > U)
	}
	return sample
}

// update all opinions, adversary might not have a particular opinion
func (sim *Sim) opinion_update(thisRound int, threshold_l, threshold_u float64) {
	U := sim.getU(threshold_l, threshold_u)
	// TODO only reserve the space CurrentAdvOpinionForNode dependent on the strategy
	advMeasure := &AdvMeasure{CurrentAdvOpinionForNode: make([]bool, sim.p.n_honest), CurrentNHonestForNode: make([]int, sim.p.n_honest)}
	sample := make([][]int, sim.p.n_honest)

	k := sim.p.k
	// change k if the first round is different
	if sim.p.enableQueryk0 {
		if thisRound == 1 {
			k = sim.p.k0
		}
	}

	// - - - - - - - - - - - - - - - - - - - - - - - -
	// non-finalized nodes' query of honest nodes
	for id := 0; id < sim.p.n_honest; id++ {
		if !sim.node[id].decided {
			if (len(sim.node[id].etaHonest) != len(sim.node[id].pHonest)) || (len(sim.node[id].etaHonest) != len(sim.node[id].opinion)) {
				fmt.Println("This should not happen.")
				wait()
			}
			// get query sample for this node
			sample[id] = sim.getSample(id, k) // this may be cause for slow code but ok for now
			// provide the number of honest nodes per honest node to the adversary
			// TODO this is already stored in pHonest, so this can be removed once the code takes pHonest rather than CurrentNHonestForNode
			countHonest := 0
			for i := 0; i < k; i++ {
				if sample[id][i] < sim.p.n_honest {
					countHonest++
				}
			}
			advMeasure.CurrentNHonestForNode[id] = countHonest
			// get votes from honest Nodes
			etaHonest, pHonest := sim.querySampleHonest(sample[id], thisRound-1, id, k)
			// remember the current etaHonest of the node
			sim.node[id].etaHonest = append(sim.node[id].etaHonest, etaHonest)
			// remember the current proportion of honest nodes of the node
			sim.node[id].pHonest = append(sim.node[id].pHonest, pHonest)
		}
	}

	// - - - - - - - - - - - - - - - - - - - - - - - -
	// take measurements, adversary may decide opinion for particular honest nodes in here
	sim.makeMeasurements(advMeasure, thisRound)

	//TODO
	// sim.makeAdvDecision

	// - - - - - - - - - - - - - - - - - - - - - - - -
	// non-finalized nodes' query of adversary nodes
	// adversary nodes in this simulation are delayed compared to the honest nodes.
	for id := 0; id < sim.p.n_honest; id++ {
		if !sim.node[id].decided {
			// get votes from adversary
			etaAdversary := sim.querySampleAdversary(sim.p, sample[id], thisRound-1, id, advMeasure)
			// calculate combined eta from honest and adversary nodes
			eta := sim.node[id].etaHonest[thisRound]*sim.node[id].pHonest[thisRound] + etaAdversary*(1-sim.node[id].pHonest[thisRound])
			// remember the current eta of the node
			sim.node[id].eta = append(sim.node[id].eta, eta)
			// new opinion of the node
			sim.node[id].opinion = append(sim.node[id].opinion, eta > U)
			if thisRound != len(sim.node[id].opinion)-1 {
				fmt.Println("This should not happen. ")
				fmt.Println("round!= len(sim.node[i].opinion)-1")
				wait()
			}
		}
	}
}
