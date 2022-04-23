package main

import (
	"log"
	"time"

	"github.com/bh90210/mlsic/sieve"
	"github.com/go-audio/audio"
	"github.com/gordonklaus/portaudio"
	// "gonum.org/v1/gonum/dsp/window"
)

func main() {
	// mlsic.Time()
	// sieve.HH()
	// os.Exit(0)
	bufferSize := 512
	buf := &audio.PCMBuffer{
		Format:   audio.FormatStereo44100,
		DataType: audio.DataTypeF64,
		F64:      make([]float64, bufferSize),
	}
	// buf := &audio.FloatBuffer{
	// 	Data:   make([]float64, bufferSize),
	// 	Format: audio.FormatStereo44100,
	// 	// Format: &audio.Format{
	// 	// 	NumChannels: 2,
	// 	// 	SampleRate:  44100,
	// 	// },
	// }

	// j := &audio.PCMBuffer{}
	// // j.

	// Audio output
	portaudio.Initialize()
	defer portaudio.Terminate()
	out := make([]float32, bufferSize)
	stream, err := portaudio.OpenDefaultStream(0, 2, 44100, len(out), &out)
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Close()

	// stream.

	if err := stream.Start(); err != nil {
		log.Fatal(err)
	}
	defer stream.Stop()

	time.Sleep(1 * time.Second)

	// currentNote := 440.0
	// osc := generator.NewOsc(generator.WaveSine, currentNote, buf.Format.SampleRate)
	s := sieve.NewSieve()
	// osc.Amplitude = 0.5
	// osc.SetAttackInMs(100)

	for {

		// populate the out buffer
		if err := s.Fill(buf); err != nil {
			log.Printf("error filling up the buffer")
		}

		// buf.F64 = window.Blackman(buf.F64)
		// bh := window.BartlettHann(buf.F64)
		// for i := range buf.F64 {
		// 	buf.F64[i] = buf.F64[i] * bh[i]
		// }

		// buf.F64 = window.Lanczos(buf.F64)

		f64ToF32Copy(out, buf.F64)

		// write to the stream
		if err := stream.Write(); err != nil {
			log.Printf("error writing to stream : %v\n", err)
		}
	}
}

// portaudio doesn't support float64 so we need to copy our data over to the
// destination buffer.
func f64ToF32Copy(dst []float32, src []float64) {
	for i := range src {
		dst[i] = float32(src[i])
	}
}
