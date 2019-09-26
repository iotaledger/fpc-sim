package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/schollz/progressbar"
)

var nCORES = runtime.NumCPU()/2 - 1 // number of available cores

// select the method for how to provide parameters
func main() {
	createDirIfNotExist("data")
	removeResultFiles()
	time.Sleep(1 * time.Second)
	simulateParameterFromFile() // read parameters from input.txt
	// other methods became redundant and were archived in main_oldmethods.go
}

// - - - - methods for providing parameters and initiating the simulation - - -

// read parameters from input.txt
func simulateParameterFromFile() {
	fmt.Println(banner()) // print simulation's logo

	lock := sync.Mutex{}
	// create nCORES channels for parallel simulation, each worker = one instance of simulation that can run in parallel
	workers := make(chan bool, nCORES)

	//parse input from file into slice of params
	P := parseInput("./input.txt")
	fmt.Println("Total parameters set:", len(P))
	bar := progressbar.New(len(P)) // initialize progress bar

	for _, p := range P {
		err := p.init() // init and check
		if err != nil {
			_ = bar.Add(1) // increment progress bar
			continue       // invalid parameter set, skipping
		}

		// every time a worker is free provide a new task (parameterset)
		workers <- true
		go task(*p, &lock, workers)
		_ = bar.Add(1) // increment progress bar
	}
	// wait for all the last workers to finish
	for workerDone := 0; workerDone < nCORES; workerDone++ {
		workers <- true
	}
}

func task(p Param, lock *sync.Mutex, workers chan bool) {
	// new RNG, seed dependent on system clock
	// it is ok to define a new RNG start in task, because
	// 1) the system clock is likely to have moved on (in doubt add time.Sleep(10 * time.Millisecond))
	// 2) a new task should have a different parameterset
	sim := &Sim{
		node:   make([]*Node, p.n),
		result: NewResult(),
		p:      &p,
		rand:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	sim.runsim(false, false)

	lock.Lock() // lock the writing so nothing else can write
	sim.saveResults()
	lock.Unlock() // unlock for the next worker to write

	<-workers // free the worker
}
