// Package mlsic .
package mlsic

import (
	"time"

	"github.com/rs/zerolog/log"
)

// SampleRate should always be 44.1, the module doesn't support higher or lower sampling rates.
const SampleRate = 44100

// MaxFrequency allowed.
const MaxFrequency = 18000

const (
	// OneSpeaker 1 speaker.
	OneSpeaker = 1
	// TwoSpeakers 2 speakers.
	TwoSpeakers = 2
)

const (
	// Speaker1 speaker number one.
	Speaker1 = iota
	// Speaker2 speaker number two.
	Speaker2
	// Speaker3 speaker number three.
	Speaker3
	// Speaker4 speaker number four.
	Speaker4
	// Speaker5 speaker number five.
	Speaker5
	// Speaker6 speaker number six.
	Speaker6
	// Speaker7 speaker number seven.
	Speaker7
	// Speaker8 speaker number eight.
	Speaker8
)

// Audio is a 64 bit float slice with PCM signal values from -1.0 to 1.0.
type Audio []float64

// Reader .
type Reader interface {
	Read(Audio) (int, error)
}

// Writer .
type Writer interface {
	Write(Audio) (int, error)
}

// Renderer .
type Renderer interface {
	Render(source []Audio, name string) error
}

// Scale a number.
func Scale(unscaledValue, scaledMin, scaledMax, unscaledMin, unscaledMax float64) float64 {
	return (scaledMax-scaledMin)*(unscaledValue-unscaledMin)/(unscaledMax-unscaledMin) + scaledMin
}

const (
	// SpeakersMinPanning in the Panning() system the minimum accepted value is 0.
	// This represents the furthest "left" speaker.
	SpeakersMinPanning = 0.

	// SpeakersMaxPanning in the Panning() system the maximum accepted value is 1.
	// This represents the furthest "right" speaker.
	SpeakersMaxPanning = 1.
)

// Panning returns the ratio to multiple the signal for a particular set of speaker numbers.
// It needs the total number of speakers, the number of the particular speakers to pan,
// and the wagon from where it derives the original panning setting.
//
// Panning expects that wagon will provide a float value between zero and one.
// If for example we have four speakers (noOfSpeakers = 4) Panning will divide one by four (1/4)
// and treat each subdivisions as a speaker's width. In this example speaker's one range is from 0 to 0.25
// with middle at 0.125. Speaker's two range is from 0.25 to 0.5 etc.
func Panning(noOfSpeakers, speakerNumber int, originalPanning float64) (panning float64) {
	switch noOfSpeakers {
	// Mono.
	case OneSpeaker:
		panning = 1

		// Stereo.
	case TwoSpeakers:
		switch speakerNumber {
		// Left.
		case Speaker1:
			panning = 1 - originalPanning

			// Right.
		case Speaker2:
			panning = originalPanning
		}

		// Three and more speakers.
	default:
		// Find the width of individual speaker.
		speakerWidth := 1. / float64(noOfSpeakers)
		// Find the max width value of current speaker.
		speakerMax := speakerWidth * float64((speakerNumber + 1))
		// Find the min width value of current speaker.
		speakerMin := speakerMax - speakerWidth
		// Find current speaker's mid point.
		speakerMid := speakerMin + (speakerWidth / 2)

		switch {
		// If the panning value is within the width of this speaker and
		// above or below speaker's mid.
		case originalPanning >= speakerMin &&
			originalPanning <= speakerMax:

			if originalPanning == speakerMid {
				panning = 1
			}

			if originalPanning < speakerMid {
				unscaledNumber := speakerMid - originalPanning
				panning = 1 - Scale(unscaledNumber, SpeakersMinPanning, SpeakersMaxPanning, 0., speakerWidth)
			}

			if originalPanning > speakerMid {
				unscaledNumber := originalPanning - speakerMid
				panning = 1 - Scale(unscaledNumber, SpeakersMinPanning, SpeakersMaxPanning, 0., speakerWidth)
			}

		// If panning value is above this speaker's range.
		// This implies that there is a speaker on the right.
		case originalPanning > speakerMax &&
			originalPanning < (speakerMid+speakerWidth):

			unscaledNumber := speakerMid + speakerWidth - originalPanning
			panning = Scale(unscaledNumber, SpeakersMinPanning, SpeakersMaxPanning, 0., speakerWidth)

		// If panning value is bellow this speaker's range.
		// This implies that there is a speaker on the left.
		case originalPanning < speakerMin &&
			originalPanning > (speakerMid-speakerWidth):

			unscaledNumber := speakerMid - originalPanning
			panning = 1 - Scale(unscaledNumber, SpeakersMinPanning, SpeakersMaxPanning, 0., speakerWidth)

		// First speaker panning in case the original panning is above the middle of last speaker.
		case speakerNumber == 0 &&
			originalPanning > speakerMid+(speakerWidth*float64(noOfSpeakers-1)):

			unscaledNumber := originalPanning - (speakerMid + (speakerWidth * float64(noOfSpeakers-1)))
			panning = Scale(unscaledNumber, SpeakersMinPanning, SpeakersMaxPanning, 0., speakerWidth)

		// Last speaker panning in case the original value was bellow first speaker's middle.
		case speakerNumber+1 == noOfSpeakers &&
			originalPanning < speakerWidth/2:

			unscaledNumber := (speakerWidth / 2) - originalPanning
			panning = Scale(unscaledNumber, SpeakersMinPanning, SpeakersMaxPanning, 0., speakerWidth)

		}
	}

	log.Debug().Float64("pan", panning).Int("speakers", noOfSpeakers).Int("speaker", speakerNumber+1).Msg("panning calculation")
	return
}

// Harmonics only method is Partials(). It returns a slice of the Partial structure.
// It is important that implementations of Harmonics return the slice of partials in
// acceding order starting from the second fundamental ([2nd, 3rd, 4th]...)
type Harmonics interface {
	Partials() []Partial
}

// Partial holds all necessary information to interpret a partial vis-à-vis it's fundamental.
type Partial struct {
	// Number is the number of the partial (eg. 2nd, 3rd etc.)
	Number int
	// AmplitudeFactor is the number to multiply fundamental's amplitude
	// in order to derive partial's amplitude.
	AmplitudeFactor float64
	// Start is the starting point of a partial in relation to it's fundamental in milliseconds.
	Start time.Duration
	// Duration of the partial in milliseconds.
	Duration time.Duration
}

// SignalLengthMultiplier this is a bit lame, fix it! TODO:
const SignalLengthMultiplier = 44

// DurationInSamples .
func (p Partial) DurationInSamples() int {
	return int(SignalLengthMultiplier * p.Duration.Abs().Milliseconds())
}

// StartInSamples .
func (p Partial) StartInSamples() int {
	return int(SignalLengthMultiplier * p.Start.Abs().Milliseconds())
}
