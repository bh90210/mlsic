package midi

import (
	"strings"
	"time"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/writer"
	driver "gitlab.com/gomidi/rtmididrv"
)

type Note int8

const (
	C Note = iota + 1
	Cs
	D
	Ds
	E
	F
	Fs
	G
	Gs
	A
	As
	B
	C0
	Cs0
	D0
	Ds0
	E0
	F0
	Fs0
	G0
	Gs0
	A0
	As0
	B0
	C1
	Cs1
	D1
	Ds1
	E1
	F1
	Fs1
	G1
	Gs1
	A1
	As1
	B1
	C2
	Cs2
	D2
	Ds2
	E2
	F2
	Fs2
	G2
	Gs2
	A2
	As2
	B2
	C3
	Cs3
	D3
	Ds3
	E3
	F3
	Fs3
	G3
	Gs3
	A3
	As3
	B3
	C4
	Cs4
	D4
	Ds4
	E4
	F4
	Fs4
	G4
	Gs4
	A4
	As4
	B4
	C5
	Cs5
	D5
	Ds5
	E5
	F5
	Fs5
	G5
	Gs5
	A5
	As5
	B5
	C6
	Cs6
	D6
	Ds6
	E6
	F6
	Fs6
	G6
	Gs6
	A6
	As6
	B6
	C7
	Cs7
	D7
	Ds7
	E7
	F7
	Fs7
	G7
	Gs7
	A7
	As7
	B7
	C8
	Cs8
	D8
	Ds8
	E8
	F8
	Fs8
	G8
	Gs8
	A8
	As8
	B8
	C9
	Cs9
	D9
	Ds9
	E9
	F9
	Fs9

	Df  Note = Cs
	Ef  Note = Ds
	Gf  Note = Fs
	Af  Note = Gs
	Bf  Note = As
	Df0 Note = Cs0
	Ef0 Note = Ds0
	Gf0 Note = Fs0
	Af0 Note = Gs0
	Bf0 Note = As0
	Df1 Note = Cs1
	Ef1 Note = Ds1
	Gf1 Note = Fs1
	Af1 Note = Gs1
	Bf1 Note = As1
	Df2 Note = Cs2
	Ef2 Note = Ds2
	Gf2 Note = Fs2
	Af2 Note = Gs2
	Bf2 Note = As2
	Df3 Note = Cs3
	Ef3 Note = Ds3
	Gf3 Note = Fs3
	Af3 Note = Gs3
	Bf3 Note = As3
	Df4 Note = Cs4
	Ef4 Note = Ds4
	Gf4 Note = Fs4
	Af4 Note = Gs4
	Bf4 Note = As4
	Df5 Note = Cs5
	Ef5 Note = Ds5
	Gf5 Note = Fs5
	Af5 Note = Gs5
	Bf5 Note = As5
	Df6 Note = Cs6
	Ef6 Note = Ds6
	Gf6 Note = Fs6
	Af6 Note = Gs6
	Bf6 Note = As6
	Df7 Note = Cs7
	Ef7 Note = Ds7
	Gf7 Note = Fs7
	Af7 Note = Gs7
	Bf7 Note = As7
	Df8 Note = Cs8
	Ef8 Note = Ds8
	Gf8 Note = Fs8
	Af8 Note = Gs8
	Bf8 Note = As8
	Df9 Note = Cs9
	Ef9 Note = Ds9
)

type Channel int8

type Parameter int8

type Value int8

type Duration float64

type Preset map[Parameter]Value

type Synth string

type Player struct {
	drv midi.Driver
	in  midi.In
	out midi.Out
	wr  *writer.Writer
}

func NewPlayer(device Synth) (*Player, error) {
	drv, err := driver.New()
	if err != nil {
		return nil, err
	}

	p := &Player{
		drv: drv,
	}

	ins, _ := drv.Ins()
	for _, in := range ins {
		if strings.Contains(in.String(), string(device)) {
			p.in = in
		}
	}
	outs, _ := drv.Outs()
	for _, out := range outs {
		if strings.Contains(out.String(), string(device)) {
			p.out = out
		}
	}

	err = p.in.Open()
	if err != nil {
		return nil, err
	}

	err = p.out.Open()
	if err != nil {
		return nil, err
	}

	wr := writer.New(p.out)
	p.wr = wr

	return p, nil
}

func (play *Player) Play(c Channel, n Note, v Value, d Duration) {
	play.wr.SetChannel(uint8(c))
	writer.NoteOn(play.wr, uint8(n), uint8(v))

	go func() {
		time.Sleep(time.Millisecond * time.Duration(d))
		play.wr.SetChannel(uint8(c))
		writer.NoteOff(play.wr, uint8(n))
	}()
}

func (play *Player) Preset(c Channel, p Preset) {
	for par, v := range p {
		play.CC(par, c, v)
	}
}

func (play *Player) CC(p Parameter, c Channel, v Value) {
	play.wr.SetChannel(uint8(c))
	writer.ControlChange(play.wr, uint8(p), uint8(v))
}

func (play *Player) PC(c Channel, v Value) {
	play.wr.SetChannel(uint8(c))
	writer.ProgramChange(play.wr, uint8(v))
}
