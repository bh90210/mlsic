package render

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-audio/aiff"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/stretchr/testify/assert"
)

var testCases = map[string]struct {
	wav    Wav
	signal []*audio.PCMBuffer

	aiff       Aiff
	aiffsignal []*audio.PCMBuffer
}{
	"one channel": {
		wav: Wav{
			Filepath: "one",
			Meta: &wav.Metadata{
				Software: "Mlsic",
				Engineer: "bh90210",
			},
		},
		signal: []*audio.PCMBuffer{
			{Format: &audio.Format{NumChannels: 1, SampleRate: 44100}, DataType: audio.DataTypeF32, SourceBitDepth: 32, F32: []float32{-0.5}},
		},

		aiff: Aiff{
			Filepath: "aiffone",
		},
		aiffsignal: []*audio.PCMBuffer{
			{Format: &audio.Format{NumChannels: 1, SampleRate: 44100}, DataType: audio.DataTypeF32, SourceBitDepth: 32, F32: []float32{-0.5}},
		},
	},
	"two channels": {
		wav: Wav{
			Filepath: "two",
			Meta:     (*wav.Metadata)(nil),
		},
		signal: []*audio.PCMBuffer{
			{Format: &audio.Format{NumChannels: 1, SampleRate: 44100}, DataType: audio.DataTypeF32, SourceBitDepth: 32, F32: []float32{0.5}},
			{Format: &audio.Format{NumChannels: 1, SampleRate: 44100}, DataType: audio.DataTypeF32, SourceBitDepth: 32, F32: []float32{0.9}},
		},

		aiff: Aiff{
			Filepath: "aifftwo",
		},
		aiffsignal: []*audio.PCMBuffer{
			{Format: &audio.Format{NumChannels: 1, SampleRate: 44100}, DataType: audio.DataTypeF32, SourceBitDepth: 32, F32: []float32{0.5}},
			{Format: &audio.Format{NumChannels: 1, SampleRate: 44100}, DataType: audio.DataTypeF32, SourceBitDepth: 32, F32: []float32{0.9}},
		},
	},
}

func TestWavRender(t *testing.T) {
	a := assert.New(t)

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			filePath, err := os.MkdirTemp("", tc.wav.Filepath)
			a.NoError(err)

			tc.wav.Filepath = filePath

			err = tc.wav.Render(tc.signal)
			a.NoError(err)

			for i := 0; i < len(tc.signal); i++ {
				wave := filepath.Join(filePath, fmt.Sprintf("%v.wav", i))
				f, err := os.Open(wave)
				a.NoError(err)

				dec := wav.NewDecoder(f)
				a.NoError(dec.Err())

				dec.ReadMetadata()
				a.Equal(tc.wav.Meta, dec.Metadata)

				a.NoError(dec.Rewind())
				buf, err := dec.FullPCMBuffer()
				a.NoError(err)

				a.EqualValues(tc.signal[i].SourceBitDepth, buf.SourceBitDepth)

				i32 := tc.signal[i].AsIntBuffer()
				a.NotNil(i32.Data)
				a.NotNil(buf.Data)
				a.Equal(i32.Data, buf.Data)

				a.NoError(dec.Rewind())
				a.Equal(tc.signal[i].Format, dec.Format())

				a.NoError(f.Close())
			}
		})
	}
}

func TestAiffRender(t *testing.T) {
	a := assert.New(t)

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			filePath, err := os.MkdirTemp("", tc.aiff.Filepath)
			a.NoError(err)

			tc.aiff.Filepath = filePath

			err = tc.aiff.Render(tc.aiffsignal)
			a.NoError(err)

			for i := 0; i < len(tc.aiffsignal); i++ {
				aifff := filepath.Join(filePath, fmt.Sprintf("%v.aiff", i))
				f, err := os.Open(aifff)
				a.NoError(err)

				dec := aiff.NewDecoder(f)
				a.NoError(dec.Err())
				a.True(dec.IsValidFile())

				buf, err := dec.FullPCMBuffer()
				a.NoError(err)
				a.NotNil(buf.Data)

				a.EqualValues(tc.aiffsignal[i].SourceBitDepth, buf.SourceBitDepth)
				a.Equal(dec.Format(), buf.PCMFormat())

				i32 := tc.aiffsignal[i].AsIntBuffer()
				a.NotNil(i32.Data)

				a.Equal(i32.Data, buf.Data)

				a.Equal(tc.aiffsignal[i].Format, dec.Format())

				a.NoError(dec.Rewind())
				a.NoError(f.Close())
			}
		})
	}
}
