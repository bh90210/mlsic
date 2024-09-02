// Package render holds three mlsic.Renderer implementations (Wav, Aiff and PortAudio.)
package render

import (
	"time"

	"github.com/bh90210/mlsic"
	"github.com/gordonklaus/portaudio"
)

const bufferSize int = 512

var _ mlsic.Renderer = (*PortAudio)(nil)

// PortAudio implements mlsic.Renderer and holds all
// necessary dependencies for setting up PortAudio.
type PortAudio struct {
	// Latency is part of the portaudio.StreamDeviceParameters.
	Latency time.Duration
	// OutputDevice is the device to be used.
	OutputDevice *portaudio.DeviceInfo
	// BufferSize is part of the portaudio.StreamDeviceParameters.
	BufferSize int
	Channels   int
}

// NewPortAudio will try to initialize with a portaudio.DefaultOutputDevice()
// with the default buffer size set at 512 and latency 10.
func NewPortAudio(opts ...PortAudioOption) (pa *PortAudio, err error) {
	err = portaudio.Initialize()
	if err != nil {
		return
	}

	// Set default values to the named return value pa.
	pa = &PortAudio{
		BufferSize: bufferSize,
		Latency:    10,
		Channels:   2,
	}

	for _, opt := range opts {
		opt(pa)
	}

	if pa.OutputDevice == nil {
		var defaultOutput *portaudio.DeviceInfo
		defaultOutput, err = portaudio.DefaultOutputDevice()
		if err != nil {
			return
		}

		pa.OutputDevice = defaultOutput
	}

	return
}

// Render will render for as many channels as len(pcmBuffer).
func (p *PortAudio) Render(a []mlsic.Audio, _ string) error {
	defer portaudio.Terminate()

	parameters := portaudio.StreamParameters{
		Input: portaudio.StreamDeviceParameters{
			Device: nil,
		},

		Output: portaudio.StreamDeviceParameters{
			Device:   p.OutputDevice,
			Channels: p.Channels,
			Latency:  p.Latency,
		},

		SampleRate:      float64(44100),
		FramesPerBuffer: p.BufferSize,
	}

	var buffers [][]float32
	for _, v := range a {
		data32 := make([]float32, len(v))
		f64ToF32Copy(data32, v)
		buffers = append(buffers, data32)
	}

	stream, err := portaudio.OpenStream(parameters, buffers)
	if err != nil {
		return err
	}

	defer stream.Stop()

	finish, startCounting := make(chan bool), make(chan bool, 1)
	// Calculate when to stop based on buffers length.
	go func() {
		length := len(buffers[0])
		<-startCounting
		time.Sleep(time.Duration(length/int(parameters.SampleRate)) * time.Second)
		finish <- true
	}()

	err = stream.Start()
	if err != nil {
		return err
	}

	startCounting <- true

	err = stream.Write()
	if err != nil {
		return err
	}

	<-finish

	return nil
}

// PortAudioOption if a custom type function that accepts *PortAudio
// and is used WithXXX PortAudio options functions.
type PortAudioOption func(*PortAudio)

// WithBufferSize sets PortAudio's buffer size.
func WithBufferSize(customSize int) PortAudioOption {
	return func(s *PortAudio) {
		s.BufferSize = customSize
	}
}

// WithOutputDevice set custom output device.
func WithOutputDevice(device *portaudio.DeviceInfo) PortAudioOption {
	return func(s *PortAudio) {
		s.OutputDevice = device
	}
}

// WithLatency sets PortAudio's latency.
func WithLatency(latency time.Duration) PortAudioOption {
	return func(s *PortAudio) {
		s.Latency = latency
	}
}

// portaudio doesn't support float64 so we need to copy our data over to the
// destination buffer.
func f64ToF32Copy(dst []float32, src []float64) {
	for i := range src {
		dst[i] = float32(src[i])
	}
}
