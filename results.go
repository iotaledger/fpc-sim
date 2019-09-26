package main

import (
	"fmt"
	"math"
	"sort"
)

// AgreementRateType defines the AgreementRate of the system
type AgreementRateType float64

// IntegrityRateType defines the IntegrityRate of the system
type IntegrityRateType float64

// TerminationRateType defines the TerminationRate of the system
type TerminationRateType float64

// MeanTerminationRoundType defines the MeanTerminationRound of the system
type MeanTerminationRoundType float64

// MedianTerminationRoundType defines the MedianTerminationRound of the system
type MedianTerminationRoundType float64

// TerminationRoundType defines the number of rounds necessary until FPC concludes
type TerminationRoundType map[int]int

// MeanLastRoundType defines the MeanLastRound
type MeanLastRoundType float64

// LastRoundHistoType defines the PDF of the average number of rounds considering individual nodes
type LastRoundHistoType map[int]float64

// OnesProportionType defines the proportion of 1s after FPC is terminated (i.e., the consensus)
type OnesProportionType map[float64]int

// OnesPropEvolutionType defines the honest eta's mean for each round for each tx
type OnesPropEvolutionType map[int][]float64

// EtaEvolutionType defines the histogram of honest eta's for each round
type EtaEvolutionType [][]int //round-histo-count

// Result defines the result of an FPC simulation
type Result struct {
	LastRoundHisto         LastRoundHistoType
	TerminationRound       TerminationRoundType
	AgreementRate          AgreementRateType
	IntegrityRate          IntegrityRateType
	TerminationRate        TerminationRateType
	MeanTerminationRound   MeanTerminationRoundType
	MedianTerminationRound MedianTerminationRoundType
	MeanLastRound          MeanLastRoundType
	OnesProportion         OnesProportionType
	OnesPropEvolution      OnesPropEvolutionType
	EtaEvolution           EtaEvolutionType
}

// NewResult initializes result
func NewResult() *Result {
	return &Result{
		LastRoundHisto:    make(map[int]float64),
		TerminationRound:  make(map[int]int),
		OnesProportion:    make(map[float64]int),
		OnesPropEvolution: make(map[int][]float64),
	}
}

func (sim *Sim) initResults() {
	sim.result.AgreementRate = AgreementRateType(0.)                   //  how often voting on a voting Object results in all honest nodes reaching consensus
	sim.result.IntegrityRate = IntegrityRateType(0.)                   //  how often integrity is reached
	sim.result.TerminationRate = TerminationRateType(0.)               //  how often the protocol terminates
	sim.result.MeanTerminationRound = MeanTerminationRoundType(0.)     //  how often the protocol terminates
	sim.result.MedianTerminationRound = MedianTerminationRoundType(0.) //  how often the protocol terminates
	sim.result.MeanLastRound = MeanLastRoundType(0.)
	sim.result.OnesProportion = make(map[float64]int)             // histogram of consensus
	sim.result.TerminationRound = make(map[int]int)               // histogram of final rounds for the voting protocol
	sim.result.EtaEvolution = make([][]int, sim.p.maxTermRound+1) // histogram of etas, averaged over several txs
	for i := 0; i < sim.p.maxTermRound+1; i++ {
		etaHisto := make([]int, sim.p.k+1)
		sim.result.EtaEvolution[i] = etaHisto
	}
}

//////////////////////////// evaluation   ////////////////////////////////

func (sim *Sim) evaluateEndOfProtocol(round int) {
	consensus := getConsensusOfRound(round, sim)
	// add agreement on this voting object to the agreement rate
	sim.result.AgreementRate += AgreementRateType(float64(btoi(consensus >= 1.-sim.p.etaAgreement || consensus <= sim.p.etaAgreement)))
	// add Integrity on this voting object to the integrity rate
	sim.result.IntegrityRate += IntegrityRateType(float64(btoi(((consensus >= 1.-sim.p.etaAgreement) && (sim.p.p0 >= 0.5)) || ((consensus <= sim.p.etaAgreement) && (sim.p.p0 < 0.5)))))
	// add Termination on this voting object to the termination rate
	sim.result.TerminationRate += TerminationRateType(float64(btoi(round < sim.p.maxTermRound-1)))
	// add final consensus on this voting object into a histogram
	sim.result.OnesProportion[consensus]++
	// add when the protocol concluded on the voting object
	sim.result.TerminationRound[round]++
	// evaluate etaHonest into histograms
	sim.evaluateEtaHisto()
	// save PDF and calculate communication overhead
	for _, node := range sim.node[:sim.p.n_honest] {
		sim.result.LastRoundHisto[node.TerminationRound]++
		sim.result.MeanLastRound += MeanLastRoundType(float64(node.TerminationRound) / float64(sim.p.n_honest))
	}
}

func (sim *Sim) evaluateEndOfSamples() {
	//normalize PDF
	for k, v := range sim.result.LastRoundHisto {
		sim.result.LastRoundHisto[k] = v / float64(sim.p.n_honest*sim.p.Nrun)
	}
	// adjust agreement rate
	sim.result.AgreementRate /= AgreementRateType(float64(sim.p.Nrun))
	// adjust IntegrityRate
	sim.result.IntegrityRate /= IntegrityRateType(float64(sim.p.Nrun))
	// adjust TerminationRate
	sim.result.TerminationRate /= TerminationRateType(float64(sim.p.Nrun))
	// calculate mean termination round
	for round, amount := range sim.result.TerminationRound {
		sim.result.MeanTerminationRound += MeanTerminationRoundType(float64(round) * float64(amount) / float64(sim.p.Nrun))
	}
	// calculate median termination round
	keys := getKeys(sim.result.TerminationRound)
	counter := 1
	for i := 0; counter < sim.p.Nrun; i++ {
		sim.result.MedianTerminationRound = MedianTerminationRoundType(keys[i])
		counter += sim.result.TerminationRound[keys[i]]
	}
	// adjust mean end round
	sim.result.MeanLastRound /= MeanLastRoundType(float64(sim.p.Nrun))
}

// provides histogram style consensus value
func getConsensusOfRound(round int, sim *Sim) float64 {
	consensus := 0.
	for id := 0; id < sim.p.n_honest; id++ {
		selectRound := round
		if sim.node[id].decided {
			selectRound = sim.node[id].TerminationRound
		}
		consensus += float64(btoi(sim.node[id].opinion[selectRound]))
	}
	if consensus == float64(sim.p.n_honest) || consensus == 0 {
		consensus /= float64(sim.p.n_honest) // return 0 or 1 only if exactly 0 or 1
	} else {
		consensus = math.Floor(10.*(consensus/float64(sim.p.n_honest)))/10. + 0.05 // else map onto histogram
	}
	return consensus
}

//////////////////////////// print results  ////////////////////////////////

func (sim *Sim) printResults() {
	// occurance of Final rounds
	keys := []int{}
	for k := range sim.result.TerminationRound {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	fmt.Println("\nTerminationRound")
	for _, k := range keys {
		fmt.Println(k, sim.result.TerminationRound[k])
	}

	// rate of consensus
	keys2 := []float64{}
	for k := range sim.result.OnesProportion {
		keys2 = append(keys2, k)
	}
	sort.Float64s(keys2)
	fmt.Println("\nConsensus Histogram")
	for _, k := range keys2 {
		s := fmt.Sprintf("%.3f", k)
		fmt.Println(s, sim.result.OnesProportion[k])
	}

	// Agreementrate
	fmt.Println("\nAgreementRate=", sim.result.AgreementRate)
}

func (r *Result) String() string {
	return fmt.Sprintln(
		r.LastRoundHisto,   // PDF of the average Number of rounds per node
		r.TerminationRound, // how many rounds are necessary to complete FPC
		r.OnesProportion,   // consensus results, i.e., the proportion of 1s
		r.AgreementRate,    // agreement of the system
		r.IntegrityRate,    // integrity of the system
		r.TerminationRate,  // termination of the system
	)
}

//////////////////////////// Stringer methods ////////////////////////////////

func (d TerminationRoundType) String() string {
	keys := []int{}
	for k := range d {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	output := fmt.Sprintln("\nNumber of rounds before FPC terminates")
	for _, k := range keys {
		output += fmt.Sprintf("%d\t%d\n", k, d[k])
	}
	return output
}

func (op OnesProportionType) String() string {
	keys := []float64{}
	for k := range op {
		keys = append(keys, k)
	}
	sort.Float64s(keys)
	output := fmt.Sprintln("\nOnes proportion")
	for _, k := range keys {
		output += fmt.Sprintf("%0.9f\t%d\n", k, op[k])
	}
	return output
}

func (op OnesPropEvolutionType) String() string {
	output := ""
	for tx, onesProportion := range op {
		output += fmt.Sprintln("\nOnes proportion of tx", tx)
		for round, v := range onesProportion {
			output += fmt.Sprintf("%d\t%0.9f\n", round, v)
		}
	}
	return output
}

func (pdf LastRoundHistoType) String() string {
	keys := []int{}
	for k := range pdf {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	output := fmt.Sprintln("PDF: Average number of nodes for a given number of rounds")
	for _, k := range keys {
		output += fmt.Sprintf("%d\t%0.6f\n", k, pdf[k])
	}
	return output
}

func (s AgreementRateType) String() string {
	return fmt.Sprintf("\nOverall Agreement:\t%0.9f\n", s)
}
func (s IntegrityRateType) String() string {
	return fmt.Sprintf("\nOverall Integrity:\t%0.9f\n", s)
}
func (s TerminationRateType) String() string {
	return fmt.Sprintf("\nOverall Termination:\t%0.9f\n", s)
}

//////////////////////////// End Stringer methods ////////////////////////////////

//////////////////////////// save Results ////////////////////////////////

func removeResultFiles() {
	deleteFile("data/result_all.csv") // this is a huge file if a matrix is calculated
	deleteFile("data/result_AgreementRate.csv")
	deleteFile("data/result_IntegrityRate.csv")
	deleteFile("data/result_TerminationRate.csv")
	deleteFile("data/result_TerminationRound.csv")
	deleteFile("data/result_MeanTerminationRound.csv")
	deleteFile("data/result_MedianTerminationRound.csv")
	deleteFile("data/result_MeanLastRound.csv")
	deleteFile("data/result_LastRoundHisto.csv")
	deleteFile("data/result_eta.csv")
}

func (sim *Sim) saveResults() {
	// all data
	// writeCSV(sim.csvAll(), "all",true)
	// Agreementrate
	writeCSV(sim.csvAgreementRate(), "AgreementRate", true)
	// IntegrityRate
	writeCSV(sim.csvIntegrityRate(), "IntegrityRate", true)
	// TerminationRate
	writeCSV(sim.csvTerminationRate(), "TerminationRate", true)
	// MeanTerminationRound
	writeCSV(sim.csvMeanTerminationRound(), "MeanTerminationRound", true)
	// MedianTerminationRound
	writeCSV(sim.csvMedianTerminationRound(), "MedianTerminationRound", true)
	// MeanLastRound
	writeCSV(sim.csvMeanLastRound(), "MeanLastRound", true)
	// Final round
	writeCSV(sim.csvTerminationRound(), "TerminationRound", true)
	// PDF Final rounds
	writeCSV(sim.csvLastRoundHisto(), "LastRoundHisto", true)

	// eta
	if sim.p.enableSaveEta {
		writeCSV(sim.csvEta(), "eta", false)
	}

}

//////////////////////////// end saveResults ////////////////////////////////

//////////////////////////// CSV methods //////////////////////////////////////

func getCsvHeader() []string {
	return []string{
		"a", "b", "beta", "k", "m", "l", "N", "p_0", "q", "Adv_strategy", "rateRandomness", "deltaWS", "gammaWS", "maxTermRound", "Type", "X", "Y",
	}
}

func (sim *Sim) csvAll() (output [][]string) {
	// output = append(output, [][]string{getCSVHeader()}...)
	output = append(output, sim.result.LastRoundHisto.csv(sim.p)...)
	output = append(output, sim.result.OnesProportion.csv(sim.p)...)
	output = append(output, sim.result.OnesPropEvolution.csv(sim.p)...)
	output = append(output, sim.result.TerminationRound.csv(sim.p)...)
	output = append(output, [][]string{sim.result.AgreementRate.csv(sim.p)}...)
	return output
}

// Agreement Rate
func (sim *Sim) csvAgreementRate() (output [][]string) {
	output = append(output, [][]string{sim.result.AgreementRate.csv(sim.p)}...)
	return output
}

// IntegrityRate
func (sim *Sim) csvIntegrityRate() (output [][]string) {
	output = append(output, [][]string{sim.result.IntegrityRate.csv(sim.p)}...)
	return output
}

// TerminationRate
func (sim *Sim) csvTerminationRate() (output [][]string) {
	output = append(output, [][]string{sim.result.TerminationRate.csv(sim.p)}...)
	return output
}

// MeanTerminationRond
func (sim *Sim) csvMeanTerminationRound() (output [][]string) {
	output = append(output, [][]string{sim.result.MeanTerminationRound.csv(sim.p)}...)
	return output
}

// MedianTerminationRond
func (sim *Sim) csvMedianTerminationRound() (output [][]string) {
	output = append(output, [][]string{sim.result.MedianTerminationRound.csv(sim.p)}...)
	return output
}

// MeanLastRound
func (sim *Sim) csvMeanLastRound() (output [][]string) {
	output = append(output, [][]string{sim.result.MeanLastRound.csv(sim.p)}...)
	return output
}

// Final Round
func (sim *Sim) csvTerminationRound() (output [][]string) {
	output = append(output, sim.result.TerminationRound.csv(sim.p)...)
	return output
}

// Final Round
func (sim *Sim) csvLastRoundHisto() (output [][]string) {
	output = append(output, sim.result.LastRoundHisto.csv(sim.p)...)
	return output
}

// Eta Histogram Evolution
func getCsvEtaHeader() []string {
	return []string{
		"Round", "eta...",
	}
}
func (sim *Sim) csvEta() (output [][]string) {
	output = append(output, sim.result.EtaEvolution.csv(sim.p)...)
	return output
}

//  csv methods for each type
func (d TerminationRoundType) csv(p *Param) (output [][]string) {
	keys := []int{}
	for k := range d {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		record := p.getCsvParam()
		record = append(record, []string{
			fmt.Sprintf("TerminationRound"),
			fmt.Sprintf("%d", k),
			fmt.Sprintf("%d", d[k]),
		}...)
		output = append(output, record)
	}
	return output
}

func (op OnesProportionType) csv(p *Param) (output [][]string) {
	keys := []float64{}
	for k := range op {
		keys = append(keys, k)
	}
	sort.Float64s(keys)
	for _, k := range keys {
		record := p.getCsvParam()
		record = append(record, []string{
			fmt.Sprintf("Ones"),
			fmt.Sprintf("%0.9f", k),
			fmt.Sprintf("%d", op[k]),
		}...)
		output = append(output, record)
	}
	return output
}

func (op OnesPropEvolutionType) csv(p *Param) (output [][]string) {
	for _, eachTx := range op {
		for i, v := range eachTx {
			record := p.getCsvParam()
			record = append(record, []string{
				fmt.Sprintf("Evo"),
				fmt.Sprintf("%d", i),
				fmt.Sprintf("%0.9f", v),
			}...)
			output = append(output, record)
		}
	}
	return output
}

func (op EtaEvolutionType) csv(p *Param) (output [][]string) {
	for i, eachRound := range op {
		record := []string{fmt.Sprintf("%d", i)}
		for _, v := range eachRound {
			record = append(record, fmt.Sprintf("%d", v))
		}
		output = append(output, record)
	}
	return output
}

func (pdf LastRoundHistoType) csv(p *Param) (output [][]string) {
	keys := []int{}
	for k := range pdf {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		record := p.getCsvParam()
		record = append(record, []string{
			fmt.Sprintf("PDF"),
			fmt.Sprintf("%d", k),
			fmt.Sprintf("%0.6f", pdf[k]),
		}...)
		output = append(output, record)
	}
	return output
}

func (s AgreementRateType) csv(p *Param) []string {
	record := p.getCsvParam()
	record = append(record, []string{
		fmt.Sprintf("Agreement Rate"),
		fmt.Sprintf(""),
		fmt.Sprintf("%0.9f", s),
	}...)
	return record
}
func (s IntegrityRateType) csv(p *Param) []string {
	record := p.getCsvParam()
	record = append(record, []string{
		fmt.Sprintf("Integrity Rate"),
		fmt.Sprintf(""),
		fmt.Sprintf("%0.9f", s),
	}...)
	return record
}
func (s TerminationRateType) csv(p *Param) []string {
	record := p.getCsvParam()
	record = append(record, []string{
		fmt.Sprintf("Termination Rate"),
		fmt.Sprintf(""),
		fmt.Sprintf("%0.9f", s),
	}...)
	return record
}
func (s MeanTerminationRoundType) csv(p *Param) []string {
	record := p.getCsvParam()
	record = append(record, []string{
		fmt.Sprintf("MeanTerminationRound"),
		fmt.Sprintf(""),
		fmt.Sprintf("%0.9f", s),
	}...)
	return record
}
func (s MedianTerminationRoundType) csv(p *Param) []string {
	record := p.getCsvParam()
	record = append(record, []string{
		fmt.Sprintf("MedianTerminationRound"),
		fmt.Sprintf(""),
		fmt.Sprintf("%0.9f", s),
	}...)
	return record
}
func (s MeanLastRoundType) csv(p *Param) []string {
	record := p.getCsvParam()
	record = append(record, []string{
		fmt.Sprintf("MeanLastRound"),
		fmt.Sprintf(""),
		fmt.Sprintf("%0.9f", s),
	}...)
	return record
}

//////////////////////////// End CSV methods //////////////////////////////////////
