package mlsic

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSignal(t *testing.T) {
	want := Audio{0, 0.06264832417874369, 0.1250505236945281, 0.18696144082725336, 0.2481378479437379, 0.3083394030591003, 0.3673295940613788, 0.42487666788983836, 0.4807545410165316, 0.5347436876541295, 0.5866320022005455, 0.6362156325320928, 0.6832997808714386, 0.7276994690840092, 0.7692402653962486, 0.8077589696806923, 0.8431042546155973, 0.8751372602002242, 0.9037321392901134, 0.9287765520091285, 0.9501721070958866, 0.9678347484506663, 0.9816950853641806, 0.991698665130853, 0.9978061869755915, 0.9999936564536083, 0.9982524797167027, 0.9925894972756649, 0.9830269571261583, 0.969602427343661, 0.9523686484908523, 0.9313933264172864, 0.9067588662653752, 0.878562048727684, 0.8469136498274179, 0.811938005715857, 0.773772524196507, 0.7325671448950292, 0.688483750195759, 0.6416955292590559, 0.5923862976180484, 0.5407497750278535, 0.48698882440436786, 0.43131465484258347, 0.3739459918455135, 0.31510821802362127, 0.2550324876406651, 0.19395481848461218, 0.13211516463136846, 0.06975647374412641, 0.007123732611892247, -0.055536995323039316, -0.11797953666742966, -0.17995857521140174, -0.24123061569515172, -0.30155494042174913, -0.360694554957803, -0.41841711920665464, -0.4744958601962072, -0.5287104629953566, -0.5808479362589023, -0.6307034490004961, -0.6780811353062289, -0.7227948638273908, -0.7646689690293299, -0.8035389413235665, -0.8392520733718676, -0.8716680600231558, -0.9006595495263083, -0.9261126438533072, -0.9479273461671313, -0.9660179536764467, -0.980313394333687, -0.9907575060537568, -0.9973092573563939, -0.9999429085653609, -0.9986481129311697, -0.9934299572800558, -0.9843089420295122, -0.9713209006488905, -0.9545168588814844, -0.9339628342811638, -0.9097395768511175, -0.8819422518036468, -0.8506800656873403, -0.8160758373504614, -0.7782655154260832, -0.7373976442346158, -0.6936327802020255, -0.6471428610864344, -0.598110530491217, -0.546728420318345, -0.49319839398099236, -0.43773075334857514, -0.3805434125398613, -0.32186104181006747, -0.26191418489530965, -0.20093835328207885, -0.13917310096006763, -0.07686108329331846, -0.014247103707104703}

	sine := Sine{
		Frequency: 440.,
		Duration:  time.Duration(3 * time.Millisecond),
	}

	i, got := sine.Signal()

	assert.Equal(t, 101, i)
	assert.Equal(t, want, got)
}

func TestDurationInSamples(t *testing.T) {
	sine := Sine{
		Duration: time.Duration(1000 * time.Millisecond),
	}

	got := sine.DurationInSamples()

	assert.Equal(t, 44000, got)
}

func TestScale(t *testing.T) {
	got := Scale(10., 0., 1., 0., 100.)
	assert.Equal(t, 0.1, got)

	got = Scale(.5, 0., 100., 0., 1.)
	assert.Equal(t, 50., got)
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
