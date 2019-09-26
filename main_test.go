package main

import (
	"math/rand"
	"testing"
	"time"
)

func TestCheckAllFinished(t *testing.T) {

	// set simulation parameters
	p := &Param{
		Nrun:          1,
		n:             2,
		k:             1,
		l:             2, // l needs to be large enough
		m:             0, // cooling off period
		MaxFinalRound: 1000,
		a:             0.75,
		b:             0.85,
		beta:          0.3,
		q:             0.0,
		p0:            1.,
	}
	p.n_honest = int(float64(p.n) * (1 - p.q))
	p.n_adv = p.n - p.n_honest

	sim := &Sim{
		node: make([]*Node, p.n),
		p:    p,
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	sim.rand = r

	// Initiate slices
	for i := 0; i < p.n; i++ {
		node := &Node{
			opinion: make([]bool, 0),
		}
		sim.node[i] = node
	}

	// ## opinions at the beginning, adversary puts its opinion to 0
	for i := 0; i < p.n_honest; i++ {
		if i < int(float64(p.n_honest)*p.p0) {
			sim.node[i].opinion = append(sim.node[i].opinion, true)
		} else {
			sim.node[i].opinion = append(sim.node[i].opinion, false)
		}
	}
	for i := p.n_honest; i < p.n; i++ {
		sim.node[i].opinion = append(sim.node[i].opinion, false)
	}

	for round := 0; round < p.m+p.l; round++ {
		for i := 0; i < p.n_honest; i++ {
			if i < int(float64(p.n_honest)*p.p0) {
				sim.node[i].opinion = append(sim.node[i].opinion, true)
			} else {
				sim.node[i].opinion = append(sim.node[i].opinion, false)
			}
		}
	}

	type testInput struct {
		round    int
		expected bool
	}
	var tests = []testInput{
		{2, true},
	}

	for _, test := range tests {

		result := sim.checkAllFinished(test.round)

		if result != test.expected {
			t.Error("Should return", test.expected, "got", result, "- given input", test)
		}
	}
}
