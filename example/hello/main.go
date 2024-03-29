// This is simple example demonstrating the use of mlsic packages with a CLI.
package main

import (
	"log"
	"os"

	"github.com/bh90210/mlsic"
	"github.com/bh90210/mlsic/generator"
	"github.com/bh90210/mlsic/pan"
	"github.com/bh90210/mlsic/render"
	"github.com/go-audio/wav"
)

func main() {
	g := &generator.Det{
		TotalGraphs: 300,
		Print:       false,
		Seed:        1,
		MaxNodes:    50,
		MaxEdges:    10,
		MaxWeight:   5,
	}

	r, err := render.NewPortAudio()
	if err != nil {
		log.Fatalf("error starting PortAudio: %s", err)
	}

	// p := &pan.LinearStereo{}
	p := &pan.Mono{}

	filePath := os.Args[1:]
	if filePath == nil {
		log.Fatal("specify filepath to save the .wav rendered file")
	}

	wave := &render.Wav{
		Filepath: filePath[0],
		Meta: &wav.Metadata{
			Software: "Mlsic",
			Comments: "Computer music <3",
		},
	}

	// TODO: fix bug that wave renderer read's destructively affecting latter renderers.
	a1 := mlsic.NewAlgo1(g, r, mlsic.Algo1WithPan(p),
		mlsic.Algo1WithAdditionalRenderer(wave), mlsic.Algo1WithLogging())

	if err := a1.Run(); err != nil {
		log.Fatalf("error while running Algo1: %s", err)
	}
}
