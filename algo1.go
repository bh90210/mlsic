// Package mlsic provides compositional algorithms
// and abstractions for producing audio.
package mlsic

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	calc "github.com/bh90210/mlsic/v1/calculate"
	"github.com/go-audio/audio"
	"github.com/inconshreveable/log15"
	"golang.org/x/sync/errgroup"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

const (
	frequencyMin = 30.    // 30 Hz
	frequencyMax = 15000. // 15 kHz
	timeMin      = 12.    // 8 milliseconds.
	timeMax      = 200.   // 5 seconds.
	panMin       = 0.
	panMax       = 1.
)

// Algo1 holds all necessary dependencies for the algorithm to run.
type Algo1 struct {
	// Model is used as Model.Dump() to retrieve to total of graphs Algo1 should process.
	Graphs Graph
	// Outputs holds output renders.
	Outputs []Renderer
	// Panning is used as Panning.Apply(sine, value) to pan the sine waves generated.
	Panning Pan
	// Format (audio.Format) holds number of channels and sample-rate in Hz.
	Format *audio.Format
	// Logging enables logging during Algo1.Run().
	Logging bool

	channels int
	graphs   []*simple.WeightedUndirectedGraph
	mu       sync.Mutex
	events   map[int][]*event
}

type event struct {
	id int64

	calc.SineOptions
	pan    float32
	signal []*audio.PCMBuffer

	ga      *calc.NetworkAnalysis
	bronKer [][]graph.Node
}

// NewAlgo1 creates an Algo1 with the default values.
// It accepts options for custom setup, ie different audio format
// or multiple renderers (ie, save as WAV and play via PortAudio.)
func NewAlgo1(g Graph, r Renderer, opts ...Algo1Option) (a1 *Algo1) {
	// Set default values to the named return value a1.
	a1 = &Algo1{
		Graphs:  g,
		Outputs: []Renderer{r},

		Format: audio.FormatMono44100,

		events: make(map[int][]*event),
	}

	// Range through options and apply them.
	for _, opt := range opts {
		opt(a1)
	}

	return
}

// Run executes the Algo1.
func (a *Algo1) Run() error {
	l := log15.New("MLSIC", "Algo1")
	if !a.Logging {
		l.SetHandler(log15.DiscardHandler())
	}

	l.Info("starting Model.Dump()")

	// Get graphs from model.
	var err error
	a.graphs, err = a.Graphs.Dump()
	if err != nil {
		return err
	}

	l.Info("finished graph dumping", "number of graphs", len(a.graphs))

	var tempBuffer = make(map[int][]*audio.PCMBuffer)
	var wg errgroup.Group

	// Construct signal/music events out of each graph concurrently.
	for i, g := range a.graphs {

		i := i
		g := g

		lg := l.New("graph", i)

		wg.Go(func() error {
			lg.Info("starting processing")

			// Make a new analysis of the graph.
			result := calc.GraphNetwork(g)
			// BronKerbosch maximal cliques will help in determining ADR times.
			bronKer := topo.BronKerbosch(g)
			// mstG := simple.NewWeightedUndirectedGraph(0, 0)
			// mst := path.Prim(mstG, g)

			// Get all nodes of the graph as slice.
			nodes := graph.NodesOf(g.Nodes())

			// Create the music events and assign them nodes' IDs and graph analysis results.
			events := a.events[i]
			for _, v := range nodes {
				events = append(events, &event{id: v.ID(), ga: result, bronKer: bronKer})
			}

			// maxDurationEvent is a helper variable to aid us determine the maximum event length
			// so we can create an equivalent long audio buffer.
			var maxDurationEvent int
			var eventWg sync.WaitGroup

			eventWg.Add(len(events))

			lg.Info("starting generation of events")

			for _, e := range events {
				go func(e *event) {
					defer eventWg.Done()

					a.processEvent(e)

					// Generate signal for event.
					sine := calc.SineGeneration(a.Format, e.SineOptions)

					// Apply panning implementation provided.
					panBuffers := a.Panning.Apply(sine, e.pan)

					// Append buffers to event's signal.
					e.signal = append(e.signal, panBuffers...)

					a.mu.Lock()
					if maxDurationEvent < len(e.signal[0].F32) {
						maxDurationEvent = len(e.signal[0].F32)
					}
					a.mu.Unlock()
				}(e)
			}

			eventWg.Wait()

			lg.Info("writing events to buffers", "number of events", len(events))

			eventBuffer := a.createTempBuffer(events, maxDurationEvent)

			a.mu.Lock()
			tempBuffer[i] = eventBuffer
			a.mu.Unlock()

			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return err
	}

	l.Info("finished writing events to buffers", "number of temp buffers", len(tempBuffer), "number of audio channels", len(tempBuffer[0]))

	// TODO: compress/limit signal before appending it to the final buffers.

	finalBuffer, duration, err := a.createFinalBuffer(tempBuffer)
	if err != nil {
		return err
	}

	l.Info("rendering", "audio duration", duration.String())

	// Run the final buffer(s) through all provided renderers.
	for _, output := range a.Outputs {
		if err := output.Render(finalBuffer); err != nil {
			return err
		}
	}

	l.Info("finished rendering", "number of renderers", len(a.Outputs))

	return nil
}

func (a *Algo1) processEvent(e *event) {
	// i := rand.Intn(5)
	i := calc.Farness

	value := e.ga.NodeValue[calc.NetworkOption(i)][e.id]
	valueMinMax := e.ga.MM[calc.NetworkOption(i)]
	rangeMinMax := calc.MinMax{Min: frequencyMin, Max: frequencyMax}

	// a.mu.Lock()
	e.Freq = float32(calc.LinearScale(value, valueMinMax, rangeMinMax))
	// a.mu.Unlock()

	ampValue := e.ga.NodeValue[calc.NetworkOption(i)][e.id]
	ampValueMinMax := e.ga.MM[calc.NetworkOption(i)]
	ampRangeMinMax := calc.MinMax{Min: 0.00001, Max: 0.009}

	// a.mu.Lock()
	e.Amp = float32(calc.LinearScale(ampValue, ampValueMinMax, ampRangeMinMax))
	// a.mu.Unlock()

	aValue := e.ga.NodeValue[calc.Farness][e.id]
	aValueMinMax := e.ga.MM[calc.Farness]
	aRangeMinMax := calc.MinMax{Min: timeMin, Max: timeMax * rand.Float64()}

	dValue := e.ga.NodeValue[calc.Betweenness][e.id]
	dValueMinMax := e.ga.MM[calc.Betweenness]
	dRangeMinMax := calc.MinMax{Min: timeMin, Max: timeMax * rand.Float64()}

	rValue := e.ga.NodeValue[calc.Betweenness][e.id]
	rValueMinMax := e.ga.MM[calc.Betweenness]
	rRangeMinMax := calc.MinMax{Min: timeMin, Max: timeMax * rand.Float64()}

	// a.mu.Lock()
	e.A = calc.LinearScale(aValue, aValueMinMax, aRangeMinMax)
	// e.A = 0
	e.D = calc.LinearScale(dValue, dValueMinMax, dRangeMinMax)
	e.R = calc.LinearScale(rValue, rValueMinMax, rRangeMinMax)
	// e.R = 0
	// a.mu.Unlock()

	panValue := e.ga.NodeValue[calc.BetweennessWeighted][e.id]
	panValueMinMax := e.ga.MM[calc.BetweennessWeighted]
	panRangeMinMax := calc.MinMax{Min: panMin, Max: panMax}

	// a.mu.Lock()
	e.pan = float32(calc.LinearScale(panValue, panValueMinMax, panRangeMinMax))
	// a.mu.Unlock()
}

func (a *Algo1) createTempBuffer(events []*event, maxDurationEvent int) []*audio.PCMBuffer {
	// Prepare buffers that will hold the combined signal of all events.
	var eventBuffer []*audio.PCMBuffer

	for i := 0; i < a.channels; i++ {
		eventBuffer = append(eventBuffer,
			&audio.PCMBuffer{
				Format:         a.Format,
				DataType:       audio.DataTypeF32,
				SourceBitDepth: 32,
				F32:            make([]float32, maxDurationEvent),
			})
	}

	// Combine signals together.
	for _, e := range events {
		for c := 0; c < a.channels; c++ {
			for i, v := range e.signal[c].F32 {
				eventBuffer[c].F32[i] += v
			}
		}
	}

	return eventBuffer
}

// TODO: make this and above function one generic func.
func (a *Algo1) createFinalBuffer(tempBuffer map[int][]*audio.PCMBuffer) ([]*audio.PCMBuffer, time.Duration, error) {
	var finalBuffer []*audio.PCMBuffer

	// Init final buffer(s).
	for i := 0; i < a.channels; i++ {
		finalBuffer = append(finalBuffer,
			&audio.PCMBuffer{
				Format:         a.Format,
				DataType:       audio.DataTypeF32,
				SourceBitDepth: 32,
			})
	}

	// Append temp buffers to final buffer(s).
	for y := 0; y < len(a.graphs); y++ {
		for i, v := range finalBuffer {
			v.F32 = append(v.F32, tempBuffer[y][i].F32...)
		}
	}

	// Calculate the duration of the audio buffer(s).
	duration, err := time.ParseDuration(fmt.Sprintf("%vms", int(float64(len(finalBuffer[0].F32))*(1000/float64(a.Format.SampleRate)))))
	if err != nil {
		return nil, 0, err
	}

	return finalBuffer, duration, nil
}

// Algo1Option if a custom type function that accepts *Algo1
// and is used WithXXX options functions.
type Algo1Option func(*Algo1)

// Algo1WithAdditionalRenderer can be used to have Algo1 render with more
// renderers than the one provided.
func Algo1WithAdditionalRenderer(r Renderer) Algo1Option {
	return func(s *Algo1) {
		s.Outputs = append(s.Outputs, r)
	}
}

// Algo1WithAudioFormat can set custom format for Algo1.
func Algo1WithAudioFormat(format *audio.Format) Algo1Option {
	return func(s *Algo1) {
		s.Format = format
	}
}

// Algo1WithPan can set custom pan for Algo1.
func Algo1WithPan(pan Pan) Algo1Option {
	return func(s *Algo1) {
		s.Panning = pan
		s.channels = pan.Channels()
	}
}

// Algo1WithLogging enables logging for Algo1.
func Algo1WithLogging() Algo1Option {
	return func(s *Algo1) {
		s.Logging = true
	}
}
