package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/schollz/progressbar"
)

type Sim struct {
	node   []*Node    // list of all nodes
	rand   *rand.Rand // random number generator
	result *Result
	p      *Param // current parameter set
}

// Node is a node participating to the FPC protocol
type Node struct {
	decided              bool      // node is decided
	TerminationRound     int       // node's final round
	opinion              []bool    // Node opinion
	eta                  []float64 // etaValue
	etaHonest            []float64 // etaValue of Honest nodes
	pHonest              []float64 // proportion of Honest nodes
	neighborsID          []int     // list of neighbors
	mana                 float64   // mana of the node
	neighborsCumuManaVec []float64 // neighbors mana vector (accumulated)
}

// main body of the simulation
func (sim *Sim) runsim(print bool, progressBar bool) {
	sim.initResults()
	bar := progressbar.New(sim.p.Nrun) // initialize progress bar

	// vote individually on Nrun vote objects
	for K := 0; K < sim.p.Nrun; K++ {

		round := 0
		sim.initiateNodes()
		onesPropEvolution := []float64{sim.p.p0} //initiate

		// ## first round
		round++
		sim.opinion_update(round, sim.p.a, sim.p.a+sim.p.deltaab)
		onesPropEvolution = append(onesPropEvolution, getConsensusOfRound(round, sim))

		// loop while at least one node is not finalized OR as long as the protocol has not reached maxTermRound
		for (!sim.checkAllFinished(round)) && (round < sim.p.maxTermRound) {
			round++
			sim.opinion_update(round, sim.p.beta, 1-sim.p.beta)
			onesPropEvolution = append(onesPropEvolution, getConsensusOfRound(round, sim))
		}

		// set unfinished nodes' final round to the current round.
		sim.finalizeUnfinishedNodes(round)

		// - - - various measurements - - -
		// add the final ones Proportion evolution
		sim.result.OnesPropEvolution[K] = onesPropEvolution
		sim.evaluateEndOfProtocol(round)

		if progressBar {
			bar.Add(1) // increment progress bar
		}
	}

	// modify some results
	sim.evaluateEndOfSamples()
	// sim.printResults()

}

// initiate nodes for a sample of the protocol
func (sim *Sim) initiateNodes() {

	// set adversary nodes random
	if sim.p.enableRandN_adv {
		sim.p.n_honest = 0
		for i := 0; i < sim.p.n; i++ {
			if sim.rand.Float64() > sim.p.q {
				sim.p.n_honest++
			}
		}
		sim.p.n_adv = sim.p.n - sim.p.n_honest
	}

	for id := 0; id < sim.p.n_honest; id++ {
		// can this be combined with the 3rd line below ?? is this realy necessary, below for eta we don't do this
		node := &Node{
			opinion: make([]bool, 0),
		}
		sim.node[id] = node
	}
	sim.createNeighborsLists()
	if sim.p.enableZipf {
		sim.distributeMana()
		sim.createProbabilityVectors()
	}
	sim.assignInitialOpinion()
}

// assign initial opinions
func (sim *Sim) assignInitialOpinion() {
	// ## opinions at the beginning, adversary puts its opinion to 0

	cutOffPoint := int(float64(sim.p.n_honest) * sim.p.p0)
	if sim.p.enableZipf {
		numNodes := 1 //always accept the first node. this may introduce an error if the first node is too large and p0 too small
		manaNow := sim.node[0].mana
		for manaNow < float64(sim.p.n_honest)/float64(sim.p.n)*sim.p.p0 { // if true accept node
			numNodes++ // add next node because last one didn't fill p0
			if numNodes < sim.p.n_honest {
				manaNow += sim.node[numNodes-1].mana // calculate new Mana level
			} else {
				manaNow = 2. //last honest node included so exit
			}
		}
		cutOffPoint = numNodes
	}
	for id := 0; id < sim.p.n_honest; id++ {
		sim.node[id].etaHonest = append(sim.node[id].etaHonest, 0.) // not defined for round 0
		sim.node[id].pHonest = append(sim.node[id].pHonest, 0.)     // not define for round 0
		if id < cutOffPoint {
			sim.node[id].opinion = append(sim.node[id].opinion, true)
		} else {
			sim.node[id].opinion = append(sim.node[id].opinion, false)
		}
	}

}

// distribute honest mana
func (sim *Sim) distributeMana() {
	honestMana := float64(sim.p.n_honest) / float64(sim.p.n)
	sumhelp := 0.
	for id := 0; id < sim.p.n_honest; id++ {
		sim.node[id].mana = math.Pow(float64(id+1), -sim.p.sZipf)
		sumhelp += sim.node[id].mana
	}
	for id := 0; id < sim.p.n_honest; id++ {
		sim.node[id].mana = sim.node[id].mana / sumhelp * honestMana
	}
	// between adversary nodes the mana is split equally.
	// More specifically the mana per malicious node is 1/n
}

//create accumulated probabilities for quorum selection
func (sim *Sim) createProbabilityVectors() {
	for id := 0; id < sim.p.n_honest; id++ {
		//create empty vectors
		a := make([]float64, len(sim.node[id].neighborsID))
		sim.node[id].neighborsCumuManaVec = a
		//calculate accumulated mana
		for id2 := 0; id2 < len(sim.node[id].neighborsID); id2++ {
			mana := 1. / float64(sim.p.n)                       // malicious mana
			if sim.node[id].neighborsID[id2] < sim.p.n_honest { //node is honest
				mana = sim.node[sim.node[id].neighborsID[id2]].mana
			}
			addMana := 0.
			if id2 > 0 {
				addMana += sim.node[id].neighborsCumuManaVec[id2-1]
			}
			sim.node[id].neighborsCumuManaVec[id2] = mana + addMana
		}
	}
}

// define neighborhood for each node
func (sim *Sim) createNeighborsLists() {
	if sim.p.enableWS {
		sim.createNeighborsListsWS()
	} else {
		a := make([]int, sim.p.n)
		for id := 0; id < sim.p.n; id++ {
			a[id] = id
		}
		for id := 0; id < sim.p.n_honest; id++ {
			sim.node[id].neighborsID = a
		}
	}
}

// create neighborhood via Watts Strogatz graph
func (sim *Sim) createNeighborsListsWS() {
	nNeighbors := minInt(sim.p.n-1-mod(sim.p.n-1, 2), maxInt(2, int((float64(sim.p.n)*sim.p.deltaWS)/2)*2)) // having this even makes things easier
	nRewire := nNeighbors / 2
	nNotNeighbors := sim.p.n - nNeighbors - 1
	if nNotNeighbors == 0 {
		panic("Watts-Strogatz graph not suitable for this simulation setting.")
	}
	if nNotNeighbors < nRewire { // don't rewire more than there are free neighbors
		nRewire = nNotNeighbors
	}
	edgeM := make([][]bool, sim.p.n)
	removeEdgeM := make([][]bool, sim.p.n)
	addEdgeM := make([][]bool, sim.p.n)

	// initiate ring network
	for id := 0; id < sim.p.n; id++ {
		row := make([]bool, sim.p.n)
		row2 := make([]bool, sim.p.n)
		removeEdgeM[id] = row2
		row3 := make([]bool, sim.p.n)
		addEdgeM[id] = row3
		for id2 := id - nNeighbors/2; id2 <= id+nNeighbors/2; id2++ {
			row[mod(id2, sim.p.n)] = true
		}
		row[id] = false
		edgeM[id] = row
	}

	// prepare rewire matrices
	rewireSet := make([]int, nNotNeighbors)
	for id := 0; id < sim.p.n; id++ {
		counter := 0
		for id2 := 0; id2 < sim.p.n; id2++ {
			if id2 != id && !edgeM[id][id2] {
				if counter >= len(rewireSet) {
					fmt.Println(id, id2, counter, len(rewireSet))
				}
				rewireSet[counter] = id2
				counter++
			}
		}
		if counter != nNotNeighbors {
			fmt.Println(nNotNeighbors)
			fmt.Println(counter)
			panic("Should not happen")
		}
		for id2 := id + 1; id2 <= id+nRewire; id2++ {
			if sim.rand.Float64() < sim.p.gammaWS {
				randID := sim.rand.Intn(nNotNeighbors)
				removeEdgeM[id][mod(id2, sim.p.n)] = true
				removeEdgeM[mod(id2, sim.p.n)][id] = true
				addEdgeM[id][rewireSet[randID]] = true
				addEdgeM[rewireSet[randID]][id] = true
				// Note: It may be worth checking wether that edge should be added again. This has some implications.
				// For example the nNeigbors/2 neighbor to the right will never be added.
				rewireSet[randID] = mod(id2, sim.p.n)
			}
		}
	}

	// rewire matrices
	for id := 0; id < sim.p.n; id++ {
		for id2 := 0; id2 < sim.p.n; id2++ {
			if removeEdgeM[id][id2] {
				edgeM[id][id2] = false
			}
			if addEdgeM[id][id2] {
				edgeM[id][id2] = true
			}
		}
	}

	// check that matrix makes sense
	for id := 0; id < sim.p.n; id++ {
		count := 0
		for id2 := 0; id2 < sim.p.n; id2++ {
			if edgeM[id][id2] {
				count++
			}
		}
		if count == 0 {
			fmt.Println("node ", id)
			panic("Node has 0 neighbors!")
		}
	}
	// for i := 0; i < len(edgeM); i++ {
	// 	fmt.Println(btoiMatrix(edgeM)[i])
	// }

	// shuffle all-nodes list
	a := make([]int, len(sim.node))
	for i := 0; i < sim.p.n; i++ {
		a[i] = i
	}
	sim.rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	// add neigbor lists to nodes
	for id := 0; id < sim.p.n; id++ {
		if a[id] < sim.p.n_honest { // currently only take care of honest nodes
			b := make([]int, 0)
			for id2 := 0; id2 < sim.p.n; id2++ {
				if edgeM[id][id2] {
					b = append(b, a[id2])
				}
			}

			if len(b) == 0 {
				fmt.Println("\n", sim.p.gammaWS)
				fmt.Println(sim.p.deltaWS)
				fmt.Println(nNeighbors, nNotNeighbors, nRewire)
				panic("b should never be zero.")
			}
			sim.node[a[id]].neighborsID = b
		}
	}
}

// check if all are finished. saves the decision for each node
func (sim *Sim) checkAllFinished(round int) bool {
	// apply k0 query
	if sim.p.enableQueryk0 {
		AllDecidedk0 := true // initially assume all nodes are finished
		if round == 1 {
			for i := 0; i < sim.p.n_honest; i++ {
				if sim.node[i].opinion[1] == sim.node[i].opinion[0] {
					if sim.node[i].eta[0] == 1 || sim.node[i].eta[0] == 0 {
						sim.node[i].decided = true
					} else {
						AllDecidedk0 = false
					}
				} else {
					AllDecidedk0 = false
				}
			}
			if AllDecidedk0 {
				return true
			}
		}
	}

	// if we are still in the cooling phase + l no check is required
	if round < sim.p.l+sim.p.m {
		return false
	}

	AllDecided := true // initially assume all nodes are finished
	for i := 0; i < sim.p.n_honest; i++ {
		if sim.node[i].decided {
			// do nothing, this node is already decided
		} else {
			lastval := sim.node[i].opinion[round]
			decided := true
			for i2 := round - 1; i2 > round-sim.p.l; i2-- {
				if lastval != sim.node[i].opinion[i2] {
					decided = false    // this node has decided this round
					AllDecided = false // if we reach this line we have proven that there is at least one unfinished node
				}
			}
			sim.node[i].decided = decided // if still true this node has finished in this round
			if decided {
				if sim.node[i].TerminationRound == 0 {
					sim.node[i].TerminationRound = round // set final round
				} else {
					fmt.Println("This should not happen.")
					wait()
				}
			}
		}
	}
	return AllDecided
}

func (sim *Sim) finalizeUnfinishedNodes(round int) {
	if round == sim.p.maxTermRound {
		for id := 0; id < sim.p.n_honest; id++ {
			if !sim.node[id].decided {
				sim.node[id].TerminationRound = round
			}
		}
	}
}
