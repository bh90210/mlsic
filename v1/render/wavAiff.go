// Package render holds various mlsic.Renderer implementations.
// It should be used along top package compositional algos.
package render

import (
	"os"
	"strconv"
	"time"

	"github.com/bh90210/mlsic/v1"
	"github.com/go-audio/audio"
	"github.com/go-audio/transforms"
	"github.com/go-audio/wav"
)

var _ mlsic.Renderer = (*Wav)(nil)

// Wav .
type Wav struct {
	Filepath string

	Meta *wav.Metadata
}

// Render .
func (w *Wav) Render(pcmBuffer []*audio.PCMBuffer) error {
	for i, buf := range pcmBuffer {
		f, err := os.Create(w.Filepath + "/" + time.Now().String() + "_" + strconv.Itoa(i) + ".wav")
		if err != nil {
			return err
		}

		format := pcmBuffer[0].Format
		wave := wav.NewEncoder(f, format.SampleRate, 32, 1, 1)
		f32 := buf.AsFloat32Buffer()

		err = transforms.PCMScaleF32(f32, 32)
		if err != nil {
			return err
		}

		i32 := f32.AsIntBuffer()

		wave.Metadata = w.Meta
		err = wave.Write(i32)
		if err != nil {
			return err
		}

		err = wave.Close()
		if err != nil {
			return err
		}

		err = f.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
