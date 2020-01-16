package main

import "fmt"

type Param struct {
	a                         float64 // parameter a of the initial threshold
	deltaab                   float64 // delta between parameter a and b of the initial threshold
	beta                      float64 // parameter beta of the subsequents thresholds
	enableExtremeBeta         bool    // the random variable switches between min and max value instead
	q                         float64 // proportion of adversaries
	n                         int     // total number of nodes
	n_honest                  int     // total number of honest nodes
	n_adv                     int     // total number of adversary nodes
	Nrun                      int     // total number of Simulation runs
	k                         int     // number of nodes to query at each round
	enableQueryk0             bool    // if true, query k0 nodes in first round and finalize if opinion equals all k0 nodes
	k0                        int     // number of nodes to query first round
	l                         int     // number of equal consecutive rounds (opinion doesn't change)
	m                         int     // cooling off period
	p0                        float64 // proportion of 1s initial honest opinions
	maxTermRound              int     // maximum number of rounds before aborting FPC
	enableRandN_adv           bool    // if true, q is used as probability to turn honest node malicious
	EnableQueryWithRepetition bool    // allow self querying and repetition in the query sample
	advStrategy               int     // Adversary startegy
	enableSaveEta             bool    // save the eta values for the parameterset
	etaAgreement              float64 // max proportion of nodes ignored to achieve agreement
	rateRandomness            float64 // probability at which a round has a random number, otherwise threshold = mean value
	enableWS                  bool    // turn on Watts-Strogatz graph instead of complete graph
	deltaWS                   float64 // proportion of nodes that are neighbors, when setting up WS
	gammaWS                   float64 // probability to reset edges in WS graph
	enableZipf                bool    // enable Zipf-distribution for Mana
	sZipf                     float64 // Zipf parameter
}

func (p *Param) init() error {
	if err := p.Check(); err != nil {
		return err
	}
	p.n_honest = int(float64(p.n) * (1 - p.q))
	p.n_adv = p.n - p.n_honest
	return nil
}

func (p *Param) String() string {
	output := "Parameters: "
	output += fmt.Sprintf("a: %v - ", p.a)
	output += fmt.Sprintf("b: %v - ", p.a+p.deltaab)
	output += fmt.Sprintf("beta: %v - ", p.beta)

	output += fmt.Sprintf("k: %v - ", p.k)
	output += fmt.Sprintf("m: %v - ", p.m)
	output += fmt.Sprintf("l: %v - ", p.l)

	output += fmt.Sprintf("N:%v - ", p.n)
	output += fmt.Sprintf("p_0: %v - ", p.p0)
	output += fmt.Sprintf("q: %v - ", p.q)
	output += fmt.Sprintln("Adv_strategy:", p.advStrategy)

	return output
}

// Check checks that all the input parameters are ok
func (p *Param) Check() error {
	var err error
	// ??? I would argue for we do not want this to see the effects of small a
	// if p.a <= 0.5 {
	// 	return fmt.Errorf("\ncondition: 0.5 < a violated  (a = %v)", p.a)
	// }
	if p.deltaab < 0 {
		return fmt.Errorf("\ncondition: a <= b violated  (a = %v, b = %v)", p.a, p.a+p.deltaab)
	}
	if p.a+p.deltaab >= 1 {
		return fmt.Errorf("\ncondition: b < 1 violated  (b = %v)", p.a+p.deltaab)
	}
	if p.beta > 0.5 || p.beta < 0 {
		return fmt.Errorf("\ncondition: 0 <= beta <= 0.5 violated  (beta = %v)", p.beta)
	}
	if p.k <= 0 {
		return fmt.Errorf("\ncondition: k > 0 violated  (k = %v)", p.k)
	}
	if p.m < 0 {
		return fmt.Errorf("\ncondition: m >= 0 violated  (m = %v)", p.m)
	}
	if p.l <= 0 {
		return fmt.Errorf("\ncondition: l > 0 violated  (l = %v)", p.l)
	}
	if p.n <= 1 {
		return fmt.Errorf("\ncondition: n > 1 violated  (n = %v)", p.n)
	}
	//TODO: check n >= k if repetitions are not allowed
	if p.p0 < 0 || p.p0 > 1 {
		return fmt.Errorf("\ncondition: 0 <= p0 <= 1 violated  (p0 = %v)", p.p0)
	}
	if p.q < 0 || p.q >= 1 {
		return fmt.Errorf("\ncondition: 0 <= q < 1 violated  (q = %v)", p.q)
	}
	if p.advStrategy != 1 && p.advStrategy != 2 && p.advStrategy != 3 && p.advStrategy != 4 && p.advStrategy != 5 {
		return fmt.Errorf("\ngiven strategy not implemented  (strategy = %v)"+
			"\nUse -help for a list of available strategies", p.advStrategy)
	}
	if p.maxTermRound < p.m+p.l {
		return fmt.Errorf("\ncondition: maxTermRound >= m + l violated  (maxTermRound = %v)", p.maxTermRound)
	}
	return err
}
