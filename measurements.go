package main

import (
	"fmt"
	"sort"
)

// measurements
type AdvMeasure struct {
	medianEtaHonestThisRound   float64
	meanOpinionHonestLastRound float64
	CurrentAdvOpinionForNode   []bool
	CurrentNHonestForNode      []int
}

// measure the mean honest opinion of the round
func (sim *Sim) measureMeanHonestOpinionLastRound(advMeasure *AdvMeasure, round int) {
	sum := 0
	for id := 0; id < sim.p.n_honest; id++ {
		if sim.node[id].decided {
			sum += btoi(sim.node[id].opinion[sim.node[id].TerminationRound])
		} else {
			sum += btoi(sim.node[id].opinion[round])
		}
	}
	advMeasure.meanOpinionHonestLastRound = float64(sum) / float64(sim.p.n_honest)
}

// measure the median etaHonest value honest opinion of the round
func (sim *Sim) measureMedianEtaHonest(advMeasure *AdvMeasure, round int) {
	var etaHonest []float64
	for id := 0; id < sim.p.n_honest; id++ {
		if !sim.node[id].decided {
			etaHonest = append(etaHonest, sim.node[id].etaHonest[round])
		} else {
			// also include the opinion of finalized nodes into the median eta
			etaHonest = append(etaHonest, float64(btoi(sim.node[id].opinion[sim.node[id].TerminationRound])))
		}
	}
	sort.Float64s(etaHonest)
	advMeasure.medianEtaHonestThisRound = GetUnsortedMedianOfSlice(etaHonest)
}

func GetUnsortedMedianOfSlice(s []float64) float64 {
	lens := len(s)
	if lens == 0 {
		fmt.Println("what should be returned in this case.?")
		wait()
	}
	if lens == 1 {
		return s[0]
	}
	// sort.Float64s(a) // sort the numbers
	if lens%2 == 1 {
		return s[lens/2]
	} else {
		return (s[lens/2-1] + s[lens/2]) / 2
	}
}

// - - -  sortable slice with a flag - -
// s := NewSlice(1.2, 25.6, 3, 5, 4)
// sort.Sort(s)

type Slice struct {
	sort.Float64Slice
	idx  []int
	flag []bool
}

func (s Slice) Swap(i, j int) {
	s.Float64Slice.Swap(i, j)
	s.idx[i], s.idx[j] = s.idx[j], s.idx[i]
	s.flag[i], s.flag[j] = s.flag[j], s.flag[i]
}

func NewSlice(flag []bool, n ...float64) *Slice {
	s := &Slice{Float64Slice: sort.Float64Slice(n), idx: make([]int, len(n)), flag: make([]bool, len(n))}
	for i := range s.idx {
		s.idx[i] = i
		s.flag[i] = flag[i]
	}
	return s
}

// evaluate the stored etas into a histogram that averages over all txs
func (sim *Sim) evaluateEtaHisto() {
	for i := 0; i < sim.p.n_honest; i++ {
		sim.evaluateEtaHistoNode(sim.node[i])
	}
}
func (sim *Sim) evaluateEtaHistoNode(node *Node) {
	for i, v := range node.eta {
		sim.result.EtaEvolution[i][etaHistoWhichColumn(v, sim.p)]++
	}
}
func etaHistoWhichColumn(eta float64, p *Param) int {
	col := int(float64(p.k) * eta)
	return col
}
