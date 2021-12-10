package elektron

import "github.com/bh90210/mlsic/midi"

type model string

// Model
const (
	CYCLES  model = "Model:Cycles"
	SAMPLES model = "Model:Samples"
)

type voice midi.Channel

// Voices/Tracks
const (
	T1 voice = iota
	T2
	T3
	T4
	T5
	T6
)

type chords midi.Parameter

// Chords
const (
	Unisonx2 midi.Parameter = iota
	Unisonx3
	Unisonx4
	Minor
	Major
	Sus2
	Sus4
	MinorMinor7
	MajorMinor7
	MinorMajor7
	MajorMajor7
	MinorMinor7Sus4
	Dim7
	MinorAdd9
	MajorAdd9
	Minor6
	Major6
	Minorb5
	Majorb5
	MinorMinor7b5
	MajorMinor7b5
	MajorAug5
	MinorMinor7Aug5
	MajorMinor7Aug5
	Minorb6
	MinorMinor9no5
	MajorMinor9no5
	MajorAdd9b5
	MajorMajor7b5
	MajorMinor7b9no5
	Sus4Aug5b9
	Sus4AddAug5
	MajorAddb5
	Major6Add4no5
	MajorMajor76no5
	MajorMajor9no5
	Fourths
	Fifths
)

const (
	// NOTE       Parameter = 3
	TRACKLEVEL midi.Parameter = 17
	MUTE       midi.Parameter = 94
	PAN        midi.Parameter = 10
	SWEEP      midi.Parameter = 18
	CONTOUR    midi.Parameter = 19
	DELAY      midi.Parameter = 12
	REVERB     midi.Parameter = 13
	VOLUMEDIST midi.Parameter = 7
	// SWING      Parameter = 15
	// CHANCE     Parameter = 14

	// model:cycles
	MACHINE     midi.Parameter = 64
	CYCLESPITCH midi.Parameter = 65
	DECAY       midi.Parameter = 80
	COLOR       midi.Parameter = 16
	SHAPE       midi.Parameter = 17
	PUNCH       midi.Parameter = 66
	GATE        midi.Parameter = 67

	// model:samples
	PITCH        midi.Parameter = 16
	SAMPLESTART  midi.Parameter = 19
	SAMPLELENGTH midi.Parameter = 20
	CUTOFF       midi.Parameter = 74
	RESONANCE    midi.Parameter = 71
	LOOP         midi.Parameter = 17
	REVERSE      midi.Parameter = 18
)

// Reverb & Delay settings
const (
	DELAYTIME midi.Parameter = iota + 85
	DELAYFEEDBACK
	REVERBSIZE
	REVERBTONE
)

// LFO settings
const (
	LFOSPEED midi.Parameter = iota + 102
	LFOMULTIPIER
	LFOFADE
	LFODEST
	LFOWAVEFORM
	LFOSTARTPHASE
	LFORESET
	LFODEPTH
)

// Machines
const (
	KICK midi.Parameter = iota
	SNARE
	METAL
	PERC
	TONE
	CHORD
)

func PT1() midi.Preset {
	p := make(map[midi.Parameter]midi.Value)
	p[MACHINE] = midi.Value(KICK)
	p[TRACKLEVEL] = midi.Value(120)
	p[MUTE] = midi.Value(0)
	p[PAN] = midi.Value(63)
	p[SWEEP] = midi.Value(16)
	p[CONTOUR] = midi.Value(24)
	p[DELAY] = midi.Value(0)
	p[REVERB] = midi.Value(0)
	p[VOLUMEDIST] = midi.Value(60)
	p[CYCLESPITCH] = midi.Value(64)
	p[DECAY] = midi.Value(29)
	p[COLOR] = midi.Value(10)
	p[SHAPE] = midi.Value(16)
	p[PUNCH] = midi.Value(0)
	p[GATE] = midi.Value(0)
	return p
}

func PT2() midi.Preset {
	p := PT1()
	p[MACHINE] = midi.Value(SNARE)
	p[SWEEP] = midi.Value(8)
	p[CONTOUR] = midi.Value(0)
	p[DECAY] = midi.Value(40)
	p[COLOR] = midi.Value(0)
	p[SHAPE] = midi.Value(127)
	return p
}

func PT3() midi.Preset {
	p := PT1()
	p[MACHINE] = midi.Value(METAL)
	p[SWEEP] = midi.Value(48)
	p[CONTOUR] = midi.Value(0)
	p[DECAY] = midi.Value(20)
	p[COLOR] = midi.Value(16)
	p[SHAPE] = midi.Value(46)
	return p
}

func PT4() midi.Preset {
	p := PT1()
	p[MACHINE] = midi.Value(PERC)
	p[SWEEP] = midi.Value(100)
	p[CONTOUR] = midi.Value(64)
	p[DECAY] = midi.Value(26)
	p[COLOR] = midi.Value(15)
	p[SHAPE] = midi.Value(38)
	return p
}

func PT5() midi.Preset {
	p := PT1()
	p[MACHINE] = midi.Value(TONE)
	p[SWEEP] = midi.Value(38)
	p[CONTOUR] = midi.Value(52)
	p[DECAY] = midi.Value(42)
	p[COLOR] = midi.Value(22)
	p[SHAPE] = midi.Value(40)
	return p
}

func PT6() midi.Preset {
	p := PT1()
	p[MACHINE] = midi.Value(CHORD)
	p[SWEEP] = midi.Value(43)
	p[CONTOUR] = midi.Value(24)
	p[DECAY] = midi.Value(64)
	p[COLOR] = midi.Value(20)
	p[SHAPE] = midi.Value(4)
	return p
}
