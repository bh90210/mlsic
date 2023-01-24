package calculate

import (
	"math"
	"testing"

	"github.com/go-audio/audio"
	"github.com/stretchr/testify/require"
)

var linearScaleTests = map[string]struct {
	value    float64
	rMinMax  MinMax
	tMinMax  MinMax
	expected float64
	isNaN    bool
}{
	"zero": {
		value:    0.,
		rMinMax:  MinMax{0., 0.},
		tMinMax:  MinMax{0., 0.},
		expected: math.NaN(),
		isNaN:    true,
	},
	"positive value": {
		value:    50.,
		rMinMax:  MinMax{0., 100.},
		tMinMax:  MinMax{1000., 9999.},
		expected: 4999.5,
	},
	"negative value": {
		value:    -1.,
		rMinMax:  MinMax{-10., 0.},
		tMinMax:  MinMax{5., 50.},
		expected: 45.,
	},
	"reversed min/max (don't do it)": {
		value:    50.,
		rMinMax:  MinMax{100., 0.},
		tMinMax:  MinMax{9999., 1000.},
		expected: 500.,
	},
}

func TestLinearScale(t *testing.T) {
	r := require.New(t)

	for name, tc := range linearScaleTests {
		t.Run(name, func(t *testing.T) {
			result := LinearScale(tc.value, tc.rMinMax, tc.tMinMax)

			r.Equal(tc.isNaN, math.IsNaN(result))

			if !tc.isNaN {
				r.Equal(tc.expected, result)
			}
		})
	}
}

var sineGenerationTests = map[string]struct {
	opts     SineOptions
	format   *audio.Format
	expected audio.PCMBuffer
}{
	"nil format": {
		opts:     SineOptions{},
		expected: audio.PCMBuffer{},
	},
	"zero sine options": {
		opts:     SineOptions{},
		format:   &audio.Format{NumChannels: 1, SampleRate: 44100},
		expected: audio.PCMBuffer{Format: &audio.Format{NumChannels: 1, SampleRate: 44100}},
	},
	"simple": {
		opts: SineOptions{
			Freq: 440.,
			Amp:  1.,
			A:    1.,
			D:    1.,
			R:    1.,
		},
		format:   &audio.Format{NumChannels: 1, SampleRate: 44100},
		expected: audio.PCMBuffer{Format: &audio.Format{NumChannels: 1, SampleRate: 44100}, I8: []int8(nil), I16: []int16(nil), I32: []int32(nil), F32: []float32{0, 0.0014238256, 0.005684115, 0.012747372, 0.022557989, 0.03503857, 0.050090402, 0.067594014, 0.087409936, 0.1093794, 0.13332546, 0.15905392, 0.1863545, 0.21500216, 0.24475828, 0.27537242, 0.3065834, 0.33812124, 0.36970866, 0.40106264, 0.43189645, 0.46192116, 0.4908476, 0.518388, 0.54425794, 0.56817824, 0.58987653, 0.609089, 0.62556255, 0.6390562, 0.6493423, 0.6562091, 0.6594612, 0.6589217, 0.65443355, 0.64586014, 0.633087, 0.6160228, 0.59460014, 0.56877625, 0.5385338, 0.50388134, 0.46485397, 0.42151314, 0.3739472, 0.31510952, 0.25503388, 0.1939563, 0.13211673, 0.06975837, 0.0071257064, -0.055535186, -0.11797791, -0.17995712, -0.24122912, -0.3015536, -0.36069342, -0.41841617, -0.47449487, -0.52870965, -0.58084726, -0.6307029, -0.6780806, -0.72279435, -0.7646686, -0.80353856, -0.83925194, -0.87166786, -0.9006594, -0.92611265, -0.9479273, -0.9660181, -0.9803134, -0.9907575, -0.99730927, -0.9999429, -0.9986481, -0.99342984, -0.9843088, -0.9713206, -0.95451653, -0.9339625, -0.909739, -0.8819416, -0.8506794, -0.81607485, -0.7782645, -0.7373963, -0.69363123, -0.63243353, -0.5709218, -0.5094496, -0.44836032, -0.38798642, -0.32864898, -0.27065343, -0.21429105, -0.15983504, -0.107540295, -0.057643365, -0.010359201, 0.034118168, 0.07561738, 0.113989435, 0.1491089, 0.18087423, 0.2092078, 0.23405626, 0.2553904, 0.27320495, 0.28751844, 0.2983726, 0.30583197, 0.30998313, 0.310934, 0.3088127, 0.3037669, 0.29596233, 0.28558174, 0.2728236, 0.2579006, 0.24103823, 0.22247332, 0.20245226, 0.18122944, 0.15906568, 0.13622631, 0.11297952, 0.08959459, 0.066340186, 0.043482542, 0.021283729}, F64: []float64(nil), DataType: 0x0, SourceBitDepth: 0x0},
	},
}

func TestSineGeneration(t *testing.T) {
	r := require.New(t)

	for name, tc := range sineGenerationTests {
		t.Run(name, func(t *testing.T) {
			result := SineGeneration(tc.format, tc.opts)
			r.Equal(tc.expected, result)
		})
	}
}
