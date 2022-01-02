package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bh90210/mlsic/midi"
	"github.com/bh90210/mlsic/midi/elektron"
)

func partition(arr []int, low, high int) int {
	index := low - 1
	pivotElement := arr[high]
	for i := low; i < high; i++ {
		if arr[i] <= pivotElement {
			index += 1
			arr[index], arr[i] = arr[i], arr[index]
		}
	}
	arr[index+1], arr[high] = arr[high], arr[index+1]

	i++
	fmt.Println(i)

	return index + 1
}

var i uint64

// QuickSortRange Sorts the specified range within the array
func QuickSortRange(arr []int, low, high int) {
	if len(arr) <= 1 {
		return
	}

	if low < high {
		fmt.Println("1", arr)
		pivot := partition(arr, low, high)
		fmt.Println("2", arr)
		QuickSortRange(arr, low, pivot-1)
		fmt.Println("3", arr)
		QuickSortRange(arr, pivot+1, high)
		fmt.Println("4", arr)
	}
}

// QuickSort Sorts the entire array
func QuickSort(arr []int) []int {
	QuickSortRange(arr, 0, len(arr)-1)
	return arr
}

type event struct {
	channel   midi.Channel
	note      midi.Note
	value     midi.Value
	duration  midi.Duration
	preset    midi.Preset
	parameter midi.Parameter
}

func main() {
	n := 1200
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
	bass := []int{60, 55, 54, 58, 55}
	for i := 0; ; i++ {
		p.Play(elektron.T1, midi.Note(bass[i%5]), 125, 2)

		p.Play(elektron.T2, midi.Note(bass[i%5]), 125, 2)

		p.Play(elektron.T3, midi.Note(bass[i%5]), 125, 2)

		p.Play(elektron.T4, midi.Note(bass[i%5]), 125, 2)

		p.Play(elektron.T5, midi.Note(bass[i%5]), 125, 2)

		p.Play(elektron.T6, midi.Note(bass[i%5]), 125, 2)

		time.Sleep(2000 * time.Millisecond)
	}

	trigger := make(chan int, 1)
	trigger2 := make(chan int, 1)
	trigger3 := make(chan int, 1)

	timeline := make(map[int]bool)
	timeline2 := make(map[int]bool)
	timeline3 := make(map[int]bool)

	for i := 2; i <= n; i++ {
		isPrime := true
		for j := 2; j < i; j++ {
			if i%j == 0 {
				isPrime = false
			}
		}

		if isPrime {
			timeline[i] = true
		}
	}

	for i := n; i <= n*2; i++ {
		isPrime := true
		for j := 2; j < i; j++ {
			if i%j == 0 {
				isPrime = false
			}
		}

		if isPrime {
			timeline2[i%n] = true
		}
	}

	for i := n * 2; i <= n*3; i++ {
		isPrime := true
		for j := 2; j < i; j++ {
			if i%j == 0 {
				isPrime = false
			}
		}

		if isPrime {
			timeline3[i%n*2] = true
		}
	}

	// log.Fatal(timeline)
	go func() {
		for i := 1; i < n; i++ {
			if timeline[i] {
				trigger <- i
			}
			if timeline2[i] {
				trigger2 <- i
			}
			if timeline3[i] {
				trigger3 <- i
			}
			time.Sleep(60000 / time.Duration(n) * time.Millisecond)
		}
		trigger <- 0
	}()

play:
	for {
		select {
		case v := <-trigger:
			if v == 0 {
				break play
			}
			// p.Play(elektron.T1, midi.C3, 127, 0.2)
			p.Play(elektron.T1, midi.Note(v%90), 127, 0.2)
			log.Println(v, timeline[v])

		case v := <-trigger2:
			if v == 0 {
				break play
			}
			// p.Play(elektron.T2, midi.E3, 127, 0.2)
			p.Play(elektron.T2, midi.Note(v%90), 127, 0.2)
			// log.Println(v, timeline2[v])
		case v := <-trigger3:
			if v == 0 {
				break play
			}
			// p.Play(elektron.T3, midi.G3, 127, 0.2)
			p.Play(elektron.T3, midi.Note(v%90), 127, 0.2)
			// log.Println(v, timeline2[v])
		}
	}
}

func yo() {
	pause := func(t time.Duration) {
		time.Sleep(t)
	}

	p, err := midi.NewPlayer(elektron.CYCLES)
	if err != nil {
		log.Println(err)
	}

	pre := make(map[midi.Parameter]midi.Value)
	pre[elektron.PAN] = 63
	pre[elektron.SWEEP] = 66
	pre[elektron.CONTOUR] = 100
	pre[elektron.DELAY] = 30
	pre[elektron.REVERB] = 90
	pre[elektron.VOLUMEDIST] = 60
	pre[elektron.CYCLESPITCH] = 64
	pre[elektron.DECAY] = 120
	pre[elektron.COLOR] = 40
	pre[elektron.SHAPE] = 50

	p.Preset(elektron.T1, pre)

	p.CC(elektron.MACHINE, elektron.T1, elektron.KICK)
	p.CC(elektron.TRACKLEVEL, elektron.T1, 120)
	p.CC(elektron.MUTE, elektron.T1, 0)
	p.CC(elektron.PUNCH, elektron.T1, 0)
	p.CC(elektron.GATE, elektron.T1, 1)

	p.Preset(elektron.T2, pre)

	p.CC(elektron.MACHINE, elektron.T2, elektron.CHORD)
	p.CC(elektron.TRACKLEVEL, elektron.T2, 120)
	p.CC(elektron.MUTE, elektron.T2, 0)
	p.CC(elektron.PUNCH, elektron.T2, 0)
	p.CC(elektron.GATE, elektron.T2, 1)
	pause(200 * time.Millisecond)

	// // note
	// p.Play(elektron.T1, midi.C4, 120, 0.5)
	// // pause
	// pause()
	// // set CC
	// p.CC(elektron.COLOR, elektron.T1, 45)

	// p.Play(elektron.T1, midi.C4, 120, 0.5)
	// pause()
	// p.CC(elektron.COLOR, elektron.T1, 85)

	// p.Play(elektron.T1, midi.C4, 120, 0.5)
	// pause()
	// p.CC(elektron.COLOR, elektron.T1, 105)

	// p.Play(elektron.T1, midi.C4, 120, 0.5)

	s := func(arr []midi.Value) (listoflists [][]midi.Value) {
		swapped := true
		for swapped {
			swapped = false

			for i := 0; i < len(arr)-1; i++ {
				cop := make([]midi.Value, len(arr))
				copy(cop, arr)
				listoflists = append(listoflists, cop)
				if arr[i] > arr[i+1] {
					arr[i+1], arr[i] = arr[i], arr[i+1]
					swapped = true
				}
			}
		}

		return listoflists
	}

	var list []midi.Value
	list = append(list, pre[elektron.PAN])
	list = append(list, pre[elektron.SWEEP])
	list = append(list, pre[elektron.CONTOUR])
	list = append(list, pre[elektron.DELAY])
	list = append(list, pre[elektron.REVERB])
	list = append(list, pre[elektron.VOLUMEDIST])
	list = append(list, pre[elektron.CYCLESPITCH])
	list = append(list, pre[elektron.DECAY])
	list = append(list, pre[elektron.COLOR])
	list = append(list, pre[elektron.SHAPE])

	data := s(list)

	p.Play(elektron.T2, midi.A4, 120, 31000)
	p.Play(elektron.T1, midi.A2, 120, 31000)
	for k, v := range data {
		fmt.Println(k, v)

		p.CC(elektron.PAN, elektron.T1, midi.Value(v[0]))
		p.CC(elektron.SWEEP, elektron.T1, midi.Value(v[1]))
		p.CC(elektron.CONTOUR, elektron.T1, midi.Value(v[2]))
		p.CC(elektron.DELAY, elektron.T1, midi.Value(v[3]))
		p.CC(elektron.REVERB, elektron.T1, midi.Value(v[4]))
		p.CC(elektron.VOLUMEDIST, elektron.T1, midi.Value(v[5]))
		p.CC(elektron.CYCLESPITCH, elektron.T1, midi.Value(v[6]))
		p.CC(elektron.DECAY, elektron.T1, midi.Value(v[7]))
		p.CC(elektron.COLOR, elektron.T1, midi.Value(v[8]))
		p.CC(elektron.SHAPE, elektron.T1, midi.Value(v[9]))

		p.CC(elektron.PAN, elektron.T2, midi.Value(v[0]))
		p.CC(elektron.SWEEP, elektron.T2, midi.Value(v[1]))
		p.CC(elektron.CONTOUR, elektron.T2, midi.Value(v[2]))
		p.CC(elektron.DELAY, elektron.T2, midi.Value(v[3]))
		p.CC(elektron.REVERB, elektron.T2, midi.Value(v[4]))
		p.CC(elektron.VOLUMEDIST, elektron.T2, midi.Value(v[5]))
		p.CC(elektron.CYCLESPITCH, elektron.T2, midi.Value(v[6]))
		p.CC(elektron.DECAY, elektron.T2, midi.Value(v[7]))
		p.CC(elektron.COLOR, elektron.T2, midi.Value(v[8]))
		p.CC(elektron.SHAPE, elektron.T2, midi.Value(v[9]))

		// p.Play(elektron.T1, midi.C4, 120, 0.2)
		pause(500 * time.Millisecond)

	}
}
