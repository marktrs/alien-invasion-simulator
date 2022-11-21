package cmd

import (
	"errors"
	"flag"
	"log"

	"github.com/marktrs/alien-invasion-assignment/simulator"
	"github.com/marktrs/alien-invasion-assignment/world"
)

var (
	WorldFilePath   string
	TotalAliens     uint64
	TotalIterations uint64
	IsLoggerEnabled bool
)

func init() {
	flag.Uint64Var(&TotalAliens, "total-aliens", 10, "Total number of aliens for simulation")
	flag.Uint64Var(&TotalIterations, "total-iterations", 10000, "Total numbers of iteration for the simulation")
	flag.BoolVar(&IsLoggerEnabled, "logging", false, "Logging detailed alien move on iteration")
	flag.StringVar(&WorldFilePath, "world-file-path", "data/world.txt", "World map data file path contains cities location data")
	flag.Parse()
}

func Run() {
	err := validateFlagsInput()
	if err != nil {
		log.Fatal(err)
	}

	w, err := world.NewFromFile(WorldFilePath)
	if err != nil {
		log.Fatal(err)
	}

	if err = simulator.
		NewSimulator(w, &simulator.Config{
			TotalAliveAliens: TotalAliens,
			TotalIterations:  TotalIterations,
			IsLoggerEnabled:  IsLoggerEnabled,
		}).
		Start(); err != nil {
		log.Fatal(err)
	}
}

func validateFlagsInput() error {
	if WorldFilePath == "" {
		return errors.New("world file path is invalid")
	}
	if TotalAliens < 1 {
		return errors.New("total aliens should be >= 1")
	}
	if TotalIterations < 1 {
		return errors.New("total iterations should be >= 1")
	}

	return nil
}
