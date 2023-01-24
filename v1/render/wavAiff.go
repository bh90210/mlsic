// Package render holds various mlsic.Renderer implementations.
// It should be used along top package compositional algos.
package render

import (
	"fmt"
	"os"

	"github.com/bh90210/mlsic/v1"
	"github.com/go-audio/audio"
	"github.com/go-audio/transforms"
	"github.com/go-audio/wav"
)

var _ mlsic.Renderer = (*Wav)(nil)

// Wav holds relevant information for encoding and saving .wav files out of audio.PCMBuffer.
type Wav struct {
	// Filepath `/path/to/directory` where the file should be saved.
	Filepath string
	// Meta holds .wav file metadata.
	Meta *wav.Metadata
}

// Render accepts a slice of audio.PCMBuffer and creates out of each one of them a mono
// .wav file named /path/to/file/0.wav for the first channel, /path/to/file/1.wav for the second etc.
func (w *Wav) Render(pcmBuffer []*audio.PCMBuffer) error {
	for i, buf := range pcmBuffer {
		f, err := os.Create(fmt.Sprintf("%s/%v.wav", w.Filepath, i))
		if err != nil {
			return err
		}

		f32 := buf.AsFloat32Buffer()

		err = transforms.PCMScaleF32(f32, 32)
		if err != nil {
			return err
		}

		i32 := f32.AsIntBuffer()

		format := pcmBuffer[0].Format
		wave := wav.NewEncoder(f, format.SampleRate, 32, 1, 1)

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
