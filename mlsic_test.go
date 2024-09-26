package mlsic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScale(t *testing.T) {
	want := 0.1
	got := Scale(10., 0., 1., 0., 100.)
	assert.Equal(t, want, got)

	want = 50.
	got = Scale(.5, 0., 100., 0., 1.)
	assert.Equal(t, want, got)
}

func TestPanning(t *testing.T) {
	tests := map[string]struct {
		noOfSpeakers    int
		speakerNumber   int
		originalPanning float64
		want            float64
	}{
		"mono 0.":  {noOfSpeakers: OneSpeaker, speakerNumber: Speaker1, originalPanning: 0., want: 1.},
		"mono 0.5": {noOfSpeakers: OneSpeaker, speakerNumber: Speaker1, originalPanning: 0.5, want: 1.},

		"stereo, middle left":  {noOfSpeakers: TwoSpeakers, speakerNumber: Speaker1, originalPanning: 0.5, want: 0.5},
		"stereo, middle right": {noOfSpeakers: TwoSpeakers, speakerNumber: Speaker2, originalPanning: 0.5, want: 0.5},

		"stereo, hard pan left, left speaker":  {noOfSpeakers: TwoSpeakers, speakerNumber: Speaker1, originalPanning: 0., want: 1.},
		"stereo, hard pan left, right speaker": {noOfSpeakers: TwoSpeakers, speakerNumber: Speaker2, originalPanning: 0., want: 0.},

		"stereo, hard pan right, left speaker":  {noOfSpeakers: TwoSpeakers, speakerNumber: Speaker1, originalPanning: 1., want: 0.},
		"stereo, hard pan right, right speaker": {noOfSpeakers: TwoSpeakers, speakerNumber: Speaker2, originalPanning: 1., want: 1.},

		"four speakers, middle of first speaker, speaker 1": {noOfSpeakers: 4, speakerNumber: Speaker1, originalPanning: 0.125, want: 1.},
		"four speakers, middle of first speaker, speaker 2": {noOfSpeakers: 4, speakerNumber: Speaker2, originalPanning: 0.125, want: 0.},
		"four speakers, middle of first speaker, speaker 3": {noOfSpeakers: 4, speakerNumber: Speaker3, originalPanning: 0.125, want: 0.},
		"four speakers, middle of first speaker, speaker 4": {noOfSpeakers: 4, speakerNumber: Speaker4, originalPanning: 0.125, want: 0.},

		"four speakers, between first speaker and fourth speaker, speaker 1": {noOfSpeakers: 4, speakerNumber: Speaker1, originalPanning: 0., want: 0.5},
		"four speakers, between first speaker and fourth speaker, speaker 2": {noOfSpeakers: 4, speakerNumber: Speaker2, originalPanning: 0., want: 0.},
		"four speakers, between first speaker and fourth speaker, speaker 3": {noOfSpeakers: 4, speakerNumber: Speaker3, originalPanning: 0., want: 0.},
		"four speakers, between first speaker and fourth speaker, speaker 4": {noOfSpeakers: 4, speakerNumber: Speaker4, originalPanning: 0., want: 0.5},

		"four speakers, above mid on second speaker, speaker 1": {noOfSpeakers: 4, speakerNumber: Speaker1, originalPanning: 0.4, want: 0.},
		"four speakers, above mid on second speaker, speaker 2": {noOfSpeakers: 4, speakerNumber: Speaker2, originalPanning: 0.4, want: 0.8999999999999999},
		"four speakers, above mid on second speaker, speaker 3": {noOfSpeakers: 4, speakerNumber: Speaker3, originalPanning: 0.4, want: 0.10000000000000009},
		"four speakers, above mid on second speaker, speaker 4": {noOfSpeakers: 4, speakerNumber: Speaker4, originalPanning: 0.4, want: 0.},

		"four speakers, bellow mid on third speaker, speaker 1": {noOfSpeakers: 4, speakerNumber: Speaker1, originalPanning: 0.55, want: 0.},
		"four speakers, bellow mid on third speaker, speaker 2": {noOfSpeakers: 4, speakerNumber: Speaker2, originalPanning: 0.55, want: 0.2999999999999998},
		"four speakers, bellow mid on third speaker, speaker 3": {noOfSpeakers: 4, speakerNumber: Speaker3, originalPanning: 0.55, want: 0.7000000000000002},
		"four speakers, bellow mid on third speaker, speaker 4": {noOfSpeakers: 4, speakerNumber: Speaker4, originalPanning: 0.55, want: 0.},

		"four speakers, above mid on fourth speaker, speaker 1": {noOfSpeakers: 4, speakerNumber: Speaker1, originalPanning: 0.88, want: 0.020000000000000018},
		"four speakers, above mid on fourth speaker, speaker 2": {noOfSpeakers: 4, speakerNumber: Speaker2, originalPanning: 0.88, want: 0.},
		"four speakers, above mid on fourth speaker, speaker 3": {noOfSpeakers: 4, speakerNumber: Speaker3, originalPanning: 0.88, want: 0.},
		"four speakers, above mid on fourth speaker, speaker 4": {noOfSpeakers: 4, speakerNumber: Speaker4, originalPanning: 0.88, want: 0.98},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := Panning(tc.noOfSpeakers, tc.speakerNumber, tc.originalPanning)
			assert.Equal(t, tc.want, got)
		})
	}
}
