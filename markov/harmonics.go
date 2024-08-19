package markov

// Harmonics .
type Harmonics struct {
	Partials map[int]float64
}

// Naive .
func (h *Harmonics) Naive() {
	// Harmonics.
	for i := 2; i < 180; i++ {
		h.Partials[i] = float64(i) * 0.01
		if h.Partials[i] > 1. {
			h.Partials[i] -= 1.
		}
	}
}
