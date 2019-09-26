# res-sim-fpc

```
 ____  ____   ___           ____  __  _  _  
(  __)(  _ \ / __)   ___   / ___)(  )( \/ ) 
 ) _)  ) __/( (__   (___)  \___ \ )( / \/ \ 
(__)  (__)   \___)         (____/(__)\_)(_/ 
```

## Run the simulation

First, you need to [install Go](https://golang.org/doc/install) if it is not already installed on your machine. It is recommended that you use the most recent version of Go.

## Build

Clone the repository with the command:

```
git clone git@github.com:iotaledger/fpc-sim.git
```

or if you prefer https over ssh

```
git clone https://github.com/iotaledger/fpc-sim.git
```

You can build your executable (as well as cross compiling for other architectures) by running the `go build` tool inside the just cloned folder `fpc-sim`:

```
go build -o sim
```
You can then run your simulation by running:

```
./sim
```

## Parameters

Parameters are provided to the simulation via the text file input.txt


|   Parameter       |       Type    | Description    |
|-------------------|:-------------:|:--------------|    
|   nRun            |   int         | Number of runs for each simulation |
|   N               |   int         | number of nodes |
|   a               |   float64     | first threshold limits in the first round [0.5 < a ≤ b < 1] |
|   deltaab         |   float64     | difference between a and b  |
|   beta            |   float64     | threshold limits in the subsequent rounds [0 ≤ beta ≤ 0.5] |
|   p0              |   float64     | proportion of ones for the initial opinion |
|   q               |   float64     | proportion of adversaries |
|   k               |   int         | amount of queries each node sends     |
|   m               |   int         | the cooling-off period |
|   l               |   int         | the determination length  |
|   maxTermRound    |   int         | Maximum number of rounds before interrupting the protocol |
|   etaAgreement    |   float64     | Proportion of nodes ignored for agreement failure
|   rateRandomness  |   float64     | Average rate at which a random number is created. Otherwise mean threshold is used.
|   Adv_Strategy    |   int      | Adversary strategy |
| enableExtremeBeta |   bool        | Experimental: beta switches between min and max threshold |
|   enableSaveEta         |   bool       | Save Etas; this does not work for vector input |
| enableRandN_adv | bool | Turn nodes adversarial with probability q, otherwise assign by index
| enableWS | bool | Enable Watts-Strogatz graph
| deltaWS | float64 | Watts-Strogatz parameter - Proportion of network 
| gammaWS | float64 | Watts-Strogatz paramter - Rewiring probability


## Results

Currently the simulator supports the following analysis:
- *AgreementRate*: Rate at which all honest nodes conclude with the same opinion 
- *IntegrityRate*: Rate at which all honest nodes conclude with the same opinion and the final opinion is the same as the original majority
- *TerminationRate*: Rate at which all honest nodes conclude before maxTermRound
- *MeanTerminationRound*: Average round when FPC terminated
- *MedianTerminationRound*: Median round when FPC terminated
- *TerminationRound*: How many rounds are necessary to complete FPC (Histogram)
- *MeanLastRound*: Mean last round across all nodes
- *LastRoundHisto*: Histogram of nodes' individual termination rounds
- *OnesProportion*: The proportion of 1s after the protocol terminates
- *OnesPropEvolution*: Evolution of the 1s proportion
- *EtaEvolution*: Histogram of honest eta's for each round 

## Adversary's strategies

The simulator supports the following strategies:
- *Strategy 1*: cautious, always send the opposite of p0
- *Strategy 2*: cautious, send the minority of the honest nodes of the previous round

## Visualisation
Data is output in csv format in the folder data. The output data can be processed via python. The output is saved in the data folder. 

For 1 vector parameter input modify the output setting in plot.py to read the correct csv column and then run:
```
py plot.py
```

For 2 vector parameter inputs modify the output setting in plot2D.py to read the correct csv columns and then run:
```
py plot2D.py
```

For 1 paramter input (no vector) and enableSaveEta enabled, gives a heatmap of the eta's evolution:
```
py plot_eps.py
```