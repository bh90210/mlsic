package elektron

import "github.com/bh90210/mlsic/midi"

// Model
const (
	CYCLES  midi.Synth = "Model:Cycles"
	SAMPLES midi.Synth = "Model:Samples"
)

// Voices/Tracks
const (
	T1 midi.Channel = iota
	T2
	T3
	T4
	T5
	T6
)

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
	TRACKLEVEL midi.Parameter = 95
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
	KICK midi.Value = iota
	SNARE
	METAL
	PERC
	TONE
	CHORD
)

func PT1() midi.Preset {
	p := make(map[midi.Parameter]midi.Value)
	p[MACHINE] = KICK
	p[TRACKLEVEL] = 120
	p[MUTE] = 0
	p[PAN] = 63
	p[SWEEP] = 16
	p[CONTOUR] = 24
	p[DELAY] = 0
	p[REVERB] = 0
	p[VOLUMEDIST] = 60
	p[CYCLESPITCH] = 64
	p[DECAY] = 29
	p[COLOR] = 10
	p[SHAPE] = 16
	p[PUNCH] = 0
	p[GATE] = 0
	return p
}

func PT2() midi.Preset {
	p := PT1()
	p[MACHINE] = SNARE
	p[SWEEP] = 8
	p[CONTOUR] = 0
	p[DECAY] = 40
	p[COLOR] = 0
	p[SHAPE] = 127
	return p
}

func PT3() midi.Preset {
	p := PT1()
	p[MACHINE] = METAL
	p[SWEEP] = 48
	p[CONTOUR] = 0
	p[DECAY] = 20
	p[COLOR] = 16
	p[SHAPE] = 46
	return p
}

func PT4() midi.Preset {
	p := PT1()
	p[MACHINE] = PERC
	p[SWEEP] = 100
	p[CONTOUR] = 64
	p[DECAY] = 26
	p[COLOR] = 15
	p[SHAPE] = 38
	return p
}

func PT5() midi.Preset {
	p := PT1()
	p[MACHINE] = TONE
	p[SWEEP] = 38
	p[CONTOUR] = 52
	p[DECAY] = 42
	p[COLOR] = 22
	p[SHAPE] = 40
	return p
}

func PT6() midi.Preset {
	p := PT1()
	p[MACHINE] = CHORD
	p[SWEEP] = 43
	p[CONTOUR] = 24
	p[DECAY] = 64
	p[COLOR] = 20
	p[SHAPE] = 4
	return p
}
