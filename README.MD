```
 ____  ____   ___           ____  __  _  _  
(  __)(  _ \ / __)   ___   / ___)(  )( \/ ) 
 ) _)  ) __/( (__   (___)  \___ \ )( / \/ \ 
(__)  (__)   \___)         (____/(__)\_)(_/ 
```
The Fast Probabilistic Consensus simulator

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
## Literature background

[1] "Consensus in the IOTA Tangle - FPC" by Dr. Sebastian Mueller:  
https://blog.iota.org/consensus-in-the-iota-tangle-fpc-b98e0f1e8fa  
[2] "..."  
?? add link to blog post here  
[3] "FPC-BI: Fast Probabilistic Consensus within Byzantine Infrastructures" by Dr. Serguei Popov and Dr. William J. Buchanan:  
https://arxiv.org/pdf/1905.10895.pdf  
[4] "Simulations of the FPC" by Dr. Sebastian Mueller:   
http://datatreker.com/simulations-of-the-fast-probabilistic-consensus-protocol-fpc

## Introduction

The following definition is taken largely from [1]:

FPC is a leaderless probabilistic binary consensus protocol of low communicational complexity that is robust in a Byzantine infrastructure. The basic idea is majority voting. 

We assume the network to have *N* nodes indexed by *1,2,…,N* for modeling purposes. Every node *i* has an opinion or state and these opinions take value in {0,1}; being the reason we speak of a binary consensus. We denote *s_i(t)* for the opinion of the node *i* at time *t*. An honest node takes the initial opinon 1 with probability *p0*.

At every round, every node queries *k* random nodes and sets its opinion equal to the majority of the queried nodes. More particularly, for a given round each node calculates the average opinion of *k* randomly queried nodes and this value is denoted *eta*. If *eta_i(t)* of node *i* is larger than the threshold the opinion *s_i(t+1)* of that node is set 1, or 0 otherwise. The threshold is set randomly between *[a,b]* in the first round, and between *[beta,1-beta]* in the following rounds. For a given node the protocol terminates if the node had the same opinion for *l* (but at least *l+m*) consecutive rounds. Unless the round reaches round *maxTermRound*, for which the protocols halts with a termination failure. 

## Parameters

Parameters are provided to the simulation via the text file input.txt. A description and functionality of some of the parameters is provided in the blog that accompanies the release of this code at [2].


|   Parameter       |       Type    | Description    |
|-------------------|:-------------:|:--------------|    
|   nRun            |   int         | Number of runs for each simulation. Each run is a new voting object and it can be, for example, understood as a new transaction. |
|   N               |   int         | Number of nodes |
|   a               |   float64     | Lower threshold limit in the first round |
|   deltaab         |   float64     | Difference between a and and upper first threshold limit b  [0.5 < a ≤ b < 1]  |
|   beta            |   float64     | Threshold limits in the subsequent rounds [0 ≤ beta ≤ 0.5] |
|   p0              |   float64     | Proportion of nodes that have initial opinion 1, or 0 otherwise |
|   q               |   float64     | Proportion of adversaries |
|   k               |   int         | Amount of nodes each node queries     |
|   m               |   int         | The cooling-off period |
|   l               |   int         | The required consecutive rounds for a node to finalize its opinion |
|   maxTermRound    |   int         | Maximum number of rounds before terminating the protocol |
|   Adv_Strategy    |   int      | Adversary strategy |
| enableWS | bool | Enable Watts-Strogatz graph, see also [2] |
| deltaWS | float64 | Watts-Strogatz parameter - Proportion of network that can be queried |
| gammaWS | float64 | Watts-Strogatz paramter - Rewiring probability |
|   rateRandomness  |   float64     | Average rate at which a random number is created, see also [2] |
|   etaAgreement    |   float64     | Proportion of nodes ignored for agreement failure |
|   enableSaveEta         |   bool       | Save etas; this does not work if a vector input is provided |
| enableExtremeBeta |   bool        | Experimental: beta switches between min and max threshold |
| enableRandN_adv | bool | Turn nodes adversarial with probability q, otherwise assign by index |


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
- *Strategy 1*: cautious, always send the opposite of *p0*
- *Strategy 2*: cautious, send the minority of the honest nodes of the previous round

## Visualisation
Data is output in csv format in the folder data. The output data can be processed via python. The figure output is saved in the data folder in EPS format (not to be confused with the variable eta). 

For 1 vector parameter input modify the output setting in plot.py to read the correct csv column and then run:
```
py plot.py
```
The script provides currently two figures: Firstly the termination-, agreement-, and integrity-rate. And secondly the mean last round for all nodes, as well as the mean round at which the protocol terminates. Note, the total message complexity in the network can be approximated by the mean last round times the quorum size k times the number of nodes.

For 2 vector parameter inputs modify the output setting in plot2D.py to read the correct csv columns and then run:
```
py plot2D.py
```
The script provides plots with the termination-, agreement-, and integrity-rate as a function (heat map) ot the input vectors. 

For 1 paramter input (no vector) turn on enableSaveEta in the input.txt file and run: 
```
py plot_eps.py
```
In every round the number of non-finalized honest nodes are counted that have taken on a given eta value and the result is stored in a slice of histograms (*EtaEvolution*). This data is employed to output a heatmap of the eta's evolution with time. 