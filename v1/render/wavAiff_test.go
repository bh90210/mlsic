package render

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/stretchr/testify/require"
)

var tests = map[string]struct {
	wav    Wav
	signal []*audio.PCMBuffer
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
			{Format: &audio.Format{NumChannels: 1, SampleRate: 44100}, DataType: audio.DataTypeF32, SourceBitDepth: 32, F32: []float32{1}},
		},
	},
	"two channels": {
		wav: Wav{
			Filepath: "two",
			Meta:     (*wav.Metadata)(nil),
		},
		signal: []*audio.PCMBuffer{
			{Format: &audio.Format{NumChannels: 1, SampleRate: 44100}, DataType: audio.DataTypeF32, SourceBitDepth: 32, F32: []float32{0}},
			{Format: &audio.Format{NumChannels: 1, SampleRate: 44100}, DataType: audio.DataTypeF32, SourceBitDepth: 32, F32: []float32{1}},
		},
	},
}

func TestRender(t *testing.T) {
	r := require.New(t)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			filePath, err := os.MkdirTemp("", tc.wav.Filepath)
			r.NoError(err)

			tc.wav.Filepath = filePath

			err = tc.wav.Render(tc.signal)
			r.NoError(err)

			for i := 0; i < len(tc.signal); i++ {
				wave := filepath.Join(filePath, fmt.Sprintf("%v.wav", i))
				f, err := os.Open(wave)
				r.NoError(err)

				dec := wav.NewDecoder(f)
				dec.ReadMetadata()
				r.Equal(tc.wav.Meta, dec.Metadata)

				r.NoError(dec.Err())

				r.NoError(dec.Rewind())
				buf, err := dec.FullPCMBuffer()
				r.NoError(err)

				i32 := tc.signal[i].AsIntBuffer()
				r.NotNil(i32.Data)
				r.NotNil(buf.Data)
				r.Equal(i32.Data, buf.Data)

				r.NoError(dec.Rewind())
				r.Equal(tc.signal[i].Format, dec.Format())

				r.NoError(f.Close())
			}
		})
	}
}
