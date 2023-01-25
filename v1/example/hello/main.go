// Package main holds a simple example demonstrating the use of mlsic packages with a CLI.
package main

import (
	"log"
	"os"

	"github.com/bh90210/mlsic/v1"
	"github.com/bh90210/mlsic/v1/generator"
	"github.com/bh90210/mlsic/v1/pan"
	"github.com/bh90210/mlsic/v1/render"
	"github.com/go-audio/wav"
)

func main() {
	g := &generator.Det{
		TotalGraphs: 100,
		Print:       false,
		Seed:        1,
	}

	r, err := render.NewPortAudio()
	if err != nil {
		log.Fatalf("error starting PortAudio: %s", err)
	}

	p := &pan.LinearStereo{}

	filePath := os.Args[1:]
	if filePath == nil {
		log.Fatal("specify filepath to save the .wav rendered file")
	}

	wave := &render.Wav{
		Filepath: filePath[0],
		Meta: &wav.Metadata{
			Engineer: "bh90210",
			Software: "Mlsic",
			Comments: "Computer music <3",
		},
	}

	a1 := mlsic.NewAlgo1(g, r, mlsic.Algo1WithPan(p),
		mlsic.Algo1WithAdditionalRenderer(wave), mlsic.Algo1WithLogging())

	if err := a1.Run(); err != nil {
		log.Fatalf("error while running Algo1: %s", err)
	}
}
