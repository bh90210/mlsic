package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/bh90210/mlsic/audio"
	"github.com/bh90210/mlsic/infinity"
	"github.com/bh90210/mlsic/midi"
	"github.com/bh90210/mlsic/midi/elektron"
	"github.com/bh90210/mlsic/sieve"
	"github.com/bh90210/mlsic/util"
)

type tri struct {
	d float64
	n midi.Note
}

func main() {
	s := infinity.Series(5, 0)
	fmt.Println(util.Float2Int(util.Markov(util.Int2Float(s), 0, 127)))

	rec := flag.Bool("r", false, "record")
	flag.Parse()

	ready := make(chan bool, 1)

	if *rec {
		go audio.Record(ready, "out")
	} else {
		ready <- true
	}

	<-ready
	fmt.Println("Starting")

	p, err := midi.NewPlayer(elektron.CYCLES)
	if err != nil {
		log.Println(err)
	}

	p1 := elektron.PT1()
	p1[elektron.DECAY] = 105
	p1[elektron.COLOR] = 127
	p1[elektron.SHAPE] = 0
	p1[elektron.SWEEP] = 0
	p1[elektron.CONTOUR] = 99
	p1[elektron.DELAY] = 15
	p1[elektron.REVERB] = 28
	p1[elektron.GATE] = 1

	p1[elektron.LFOSPEED] = 127
	p1[elektron.LFODEPTH] = 97
	// p1[elektron.LFODEST] = 10 // set FTUN manually
	p1[elektron.LFOMULTIPIER] = 3
	p1[elektron.LFOWAVEFORM] = 0
	p.Preset(elektron.T1, p1)

	p2 := elektron.PT2()
	p2[elektron.PAN] = 20
	p2[elektron.DECAY] = 66
	p2[elektron.COLOR] = 0
	p2[elektron.SHAPE] = 0
	p2[elektron.SWEEP] = 0
	p2[elektron.CONTOUR] = 0
	p2[elektron.DELAY] = 30
	p2[elektron.REVERB] = 39
	p2[elektron.GATE] = 0
	p2[elektron.TRACKLEVEL] = 60
	p.Preset(elektron.T2, p2)

	p3 := elektron.PT3()
	p3[elektron.PAN] = 97
	p3[elektron.DECAY] = 35
	p3[elektron.COLOR] = 127
	p3[elektron.SHAPE] = 0
	p3[elektron.SWEEP] = 0
	p3[elektron.CONTOUR] = 127
	p3[elektron.DELAY] = 47
	p3[elektron.REVERB] = 30
	p3[elektron.GATE] = 0
	p3[elektron.TRACKLEVEL] = 70

	p3[elektron.LFOSPEED] = 127
	p3[elektron.LFODEPTH] = 97
	p3[elektron.LFODEST] = 13
	p3[elektron.LFOMULTIPIER] = 2
	p3[elektron.LFOWAVEFORM] = 1
	p.Preset(elektron.T3, p3)

	p4 := elektron.PT4()
	p4[elektron.PAN] = 80
	p4[elektron.PITCH] = 80 // 12
	p4[elektron.DECAY] = 26
	p4[elektron.COLOR] = 40
	p4[elektron.SHAPE] = 0
	p4[elektron.SWEEP] = 0
	p4[elektron.CONTOUR] = 10
	p4[elektron.DELAY] = 0
	p4[elektron.REVERB] = 0
	p4[elektron.GATE] = 0
	p4[elektron.TRACKLEVEL] = 85
	p.Preset(elektron.T4, p4)

	p5 := elektron.PT5()
	p5[elektron.PAN] = 55
	p5[elektron.DECAY] = 69
	p5[elektron.COLOR] = 10 // 1
	p5[elektron.SHAPE] = 32
	p5[elektron.SWEEP] = 43
	p5[elektron.CONTOUR] = 127
	p5[elektron.DELAY] = 0
	p5[elektron.REVERB] = 62
	p5[elektron.GATE] = 0
	p5[elektron.TRACKLEVEL] = 85
	p.Preset(elektron.T5, p5)

	p6 := elektron.PT6()
	p2[elektron.PAN] = 55
	p6[elektron.DECAY] = 75
	p6[elektron.COLOR] = 117
	p6[elektron.SHAPE] = 7
	p6[elektron.SWEEP] = 0
	p6[elektron.CONTOUR] = 127
	p6[elektron.DELAY] = 39
	p6[elektron.REVERB] = 104
	p6[elektron.GATE] = 0
	p6[elektron.TRACKLEVEL] = 85
	p.Preset(elektron.T6, p6)

	p.CC(elektron.DELAYFEEDBACK, elektron.T1, 42)
	p.CC(elektron.DELAYTIME, elektron.T1, 17)
	p.CC(elektron.REVERBSIZE, elektron.T1, 28)
	p.CC(elektron.REVERBTONE, elektron.T1, 27)

	triggers := make([]chan tri, 0)
	for i := 0; i < 6; i++ {
		triggers = append(triggers, make(chan tri))
	}

	for i := 0; i < 6; i++ {
		go func(track midi.Channel) {
			for {
				tri := <-triggers[track]
				p.Play(track, tri.n, 127, tri.d)
			}
		}(midi.Channel(i))
	}

	time.Sleep(500 * time.Millisecond)

	maxTriggers := 1000

	bass := infinity.Series(maxTriggers, 50)

	for k, v := range bass {
		switch {
		case v < 0:
			bass[k] = v * -1

		case v > 127:
			bass[k] = 126

		case v < 6:
			bass[k] = 6

		}
	}

	fmt.Println(bass)

	timing := sieve.Intervallic(infinity.Series(maxTriggers, maxTriggers))
	for k, v := range timing {
		timing[k] = v / maxTriggers
	}

	fmt.Println("timing", timing)
	factor := 100

	// Print total running time in minutes secods
	var sum int
	// Account for the slow down.
	for i := 0; i < 250; i++ {
		sum += i
	}
	sum = sum * 6
	for _, v := range timing {
		v = v * factor
		sum += v
	}

	// 5m0.45s
	fmt.Println(time.Duration(sum) * time.Millisecond)
	brake := 30
	for i := 0; i < maxTriggers; i++ {
		for _, c := range triggers {
			c <- tri{d: float64(timing[i] * factor), n: midi.Note(bass[i])}

			// Start shifting sounds after 750th trigger.
			if i > brake {
				time.Sleep(time.Duration((i-brake)/2) * time.Millisecond)
			}
		}

		time.Sleep(time.Duration(timing[i]*factor) * time.Millisecond)
	}

	// maxTriggers := 1000

	// bass := sieve.Intervallic(infinity.Series(maxTriggers, 25))

	// for k, v := range bass {
	// 	switch {
	// 	case v < 0:
	// 		bass[k] = v * -1

	// 	case v > 127:
	// 		bass[k] = 126

	// 	case v < 6:
	// 		bass[k] = 6

	// 	}
	// }

	// fmt.Println(bass)

	// timing := sieve.Intervallic(infinity.Series(maxTriggers, 0))

	// fmt.Println("timing", timing)
	// factor := 50
	// for i := 0; i < maxTriggers; i++ {
	// 	triggers[2] <- tri{d: float64(timing[i] * factor), n: midi.Note(bass[i])}
	// 	p3 := elektron.PT3()
	// 	// p3[elektron.PAN] = 97
	// 	// p3[elektron.DECAY] = 35
	// 	// p3[elektron.COLOR] = 127
	// 	// p3[elektron.SHAPE] = 0
	// 	// p3[elektron.SWEEP] = 0
	// 	// p3[elektron.CONTOUR] = 127
	// 	// p3[elektron.DELAY] = 47
	// 	// p3[elektron.REVERB] = 30
	// 	// p3[elektron.GATE] = 0
	// 	p3[elektron.TRACKLEVEL] = midi.Value(rand.Int31n(127))

	// 	// p3[elektron.LFOSPEED] = 127
	// 	// p3[elektron.LFODEPTH] = 97
	// 	// p3[elektron.LFODEST] = 13
	// 	// p3[elektron.LFOMULTIPIER] = 2
	// 	// p3[elektron.LFOWAVEFORM] = 1
	// 	p.Preset(elektron.T3, p3)

	// 	// triggers[3] <- tri{d: float64(timing[i] * factor), n: midi.Note(bass[i])}
	// 	// triggers[4] <- tri{d: float64(timing[i] * factor), n: midi.Note(bass[i])}

	// 	if timing[i]%2 == 0 {
	// 		triggers[5] <- tri{d: float64(timing[i] * factor), n: midi.Note(bass[i])}
	// 		// triggers[1] <- tri{d: float64(timing[i] * factor), n: midi.Note(bass[i])}
	// 	}
	// 	if timing[i]%3 == 0 {
	// 		// p2[elektron.TRACKLEVEL] = midi.Value(rand.Int31n(127))
	// 		// p.Preset(elektron.T2, p2)
	// 		p.CC(elektron.TRACKLEVEL, elektron.T2, midi.Value(rand.Int31n(127)))

	// 		// triggers[1] <- tri{d: float64(timing[i] * factor), n: midi.Note(bass[i])}

	// 	}

	// 	time.Sleep(time.Duration(timing[i]*factor) * time.Millisecond)
	// }

	//
	//
	// 3
	//
	//

	// triggers := make([]chan tri, 0)
	// for i := 0; i < 6; i++ {
	// 	triggers = append(triggers, make(chan tri))
	// }

	// for i := 0; i < 6; i++ {
	// 	go func(track midi.Channel) {
	// 		for {
	// 			tri := <-triggers[track]
	// 			p.Play(track, tri.n, 127, tri.d)
	// 		}
	// 	}(midi.Channel(i))
	// }

	// time.Sleep(500 * time.Millisecond)

	// maxTriggers := 60

	// bass := sieve.Intervallic(infinity.Series(maxTriggers, 30))

	// for k, v := range bass {
	// 	switch {
	// 	case v < 0:
	// 		bass[k] = v * -1

	// 	case v > 127:
	// 		bass[k] = 126

	// 	case v < 6:
	// 		bass[k] = 6 + (v % 126)

	// 	}
	// }

	// fmt.Println(bass)

	// timing := sieve.Intervallic(infinity.Series(maxTriggers, 0))
	// // sort.Ints(timing)
	// // fmt.Println("timing sorted", timing)
	// // os.Exit(0)
	// // fmt.Println("timing", timing)
	// factor := 500

	// // Print total running time in minutes secods
	// var sum int
	// for _, v := range timing {
	// 	v = v * factor
	// 	sum += v
	// }

	// fmt.Println(timing)
	// fmt.Println(time.Duration(sum) * time.Millisecond)

	// // for k, v := range bass {	0
	// // 	fmt.Println(v, timing[k])	// }

	// for i := 0; i < maxTriggers; i++ {
	// 	// if timing[i]%2 == 0 {
	// 	// 	triggers[3] <- tri{d: float64(timing[i] * factor), n: midi.Note(bass[i])}
	// 	// }

	// 	// if timing[i]%3 == 0 {
	// 	// 	// fmt.Println(bass[i])

	// 	// 	triggers[3] <- tri{d: float64(timing[i] * factor), n: midi.Note(bass[i])}

	// 	// 	triggers[1] <- tri{d: float64(timing[i] * factor), n: midi.Note(bass[i])}
	// 	// 	triggers[4] <- tri{d: float64(timing[i] * factor), n: midi.Note(bass[i])}

	// 	// }

	// 	// if timing[i]%4 == 0 {
	// 	// 	triggers[4] <- tri{d: float64(timing[i] * factor), n: midi.Note(bass[i])}
	// 	// }

	// 	// if timing[i]%5 == 0 {
	// 	// 	triggers[1] <- tri{d: float64(timing[i] * factor), n: midi.Note(bass[i])}
	// 	// }

	// 	// //sec
	// 	if timing[i]%2 == 0 {
	// 		triggers[3] <- tri{d: float64(timing[i] * factor), n: midi.Note(bass[i])}
	// 	}

	// 	if timing[i]%3 == 0 {
	// 		triggers[1] <- tri{d: float64(timing[i] * factor), n: midi.Note(bass[i])}
	// 		triggers[4] <- tri{d: float64(timing[i] * factor), n: midi.Note(bass[i])}

	// 	}

	// 	if timing[i]%4 == 0 {
	// 		triggers[4] <- tri{d: float64(timing[i] * factor), n: midi.Note(bass[i])}

	// 	}

	// 	if timing[i]%5 == 0 {
	// 		triggers[1] <- tri{d: float64(timing[i] * factor), n: midi.Note(bass[i])}
	// 	}

	// 	time.Sleep(time.Duration(timing[i]*factor) * time.Millisecond)
	// }
}
