package render

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/bh90210/mlsic"
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
func (w *Wav) Render(sources []mlsic.Audio) error {
	var wg sync.WaitGroup

	for i, source := range sources {
		wg.Add(1)

		go func(i int, source mlsic.Audio) {
			defer wg.Done()

			f, err := os.Create(filepath.Join(w.Filepath, fmt.Sprintf("%v.wav", i)))
			if err != nil {
				log.Fatal().Err(err).Int("file", i).Msg("creating file")
			}

			buf := audio.PCMBuffer{
				Format:         &audio.Format{NumChannels: 1, SampleRate: 44100},
				F64:            source,
				DataType:       audio.DataTypeF64,
				SourceBitDepth: 1,
			}

			f64 := buf.AsFloatBuffer()
			err = transforms.PCMScale(f64, 32)
			if err != nil {
				log.Fatal().Err(err).Int("file", i).Msg("transform")
			}

			intBuf := f64.AsIntBuffer()

			wave := wav.NewEncoder(f, buf.Format.SampleRate, 32, 1, 1)

			if w.Meta != nil {
				wave.Metadata = w.Meta
			}

			err = wave.Write(intBuf)
			if err != nil {
				log.Fatal().Err(err).Int("file", i).Msg("writing buf")
			}

			err = wave.Close()
			if err != nil {
				log.Fatal().Err(err).Int("file", i).Msg("closing buffer")
			}

			err = f.Close()
			if err != nil {
				log.Fatal().Err(err).Int("file", i).Msg("closing file")
			}

			log.Info().Int("file length", wave.WrittenBytes).Int("file", i).Msg("file length")

		}(i, source)
	}

	wg.Wait()

	return nil
}

// var _ mlsic.Renderer = (*Aiff)(nil)

// // Aiff holds relevant information for encoding and saving .aiff files out of audio.PCMBuffer.
// type Aiff struct {
// 	// Filepath `/path/to/directory` where the file should be saved.
// 	Filepath string
// }

// // Render accepts a slice of audio.PCMBuffer and creates out of each one of them a mono
// // .aiff file named /path/to/file/0.aiff for the first channel, /path/to/file/1.aiff for the second etc.
// func (a *Aiff) Render(readers []mlsic.Reader) error {
// 	for i, buf := range pcmBuffer {
// 		f, err := os.Create(filepath.Join(a.Filepath, fmt.Sprintf("%v.aiff", i)))
// 		if err != nil {
// 			return err
// 		}

// 		f32 := buf.AsFloat32Buffer()

// 		err = transforms.PCMScaleF32(f32, 32)
// 		if err != nil {
// 			return err
// 		}

// 		i32 := f32.AsIntBuffer()
// 		format := pcmBuffer[0].Format

// 		aiff := aiff.NewEncoder(f, format.SampleRate, 32, 1)

// 		err = aiff.Write(i32)
// 		if err != nil {
// 			return err
// 		}

// 		err = aiff.Close()
// 		if err != nil {
// 			return err
// 		}

// 		err = f.Close()
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
