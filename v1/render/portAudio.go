// Package render holds various mlsic.Renderer implementations.
// It should be used along top package compositional algos.
package render

import (
	"time"

	"github.com/bh90210/mlsic/v1"
	"github.com/go-audio/audio"
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
}

// NewPortAudio will try to initialize with a portaudio.DefaultOutputDevice()
// with the default buffer size set at 512 abd latency 10.
func NewPortAudio(opts ...PortAudioOption) (pa *PortAudio, err error) {
	err = portaudio.Initialize()
	if err != nil {
		return
	}

	// Set default values to the named return value pa.
	pa = &PortAudio{
		BufferSize: bufferSize,
		Latency:    10,
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
// Samplerate is set by the first channel's (pcmBuffer[0]) format.
func (p *PortAudio) Render(pcmBuffer []*audio.PCMBuffer) error {
	defer portaudio.Terminate()

	parameters := portaudio.StreamParameters{
		Input: portaudio.StreamDeviceParameters{
			Device: nil,
		},

		Output: portaudio.StreamDeviceParameters{
			Device:   p.OutputDevice,
			Channels: len(pcmBuffer),
			Latency:  p.Latency,
		},

		SampleRate:      float64(pcmBuffer[0].Format.SampleRate),
		FramesPerBuffer: p.BufferSize,
	}

	var buffers [][]float32
	for _, b := range pcmBuffer {
		buffers = append(buffers, b.F32)
	}

	stream, err := portaudio.OpenStream(parameters, buffers)
	if err != nil {
		return err
	}

	defer stream.Stop()

	// TODO: Log length, channels, sample/bitrate.

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
