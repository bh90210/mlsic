---
title: v0.1.0
weight: 1
---

# Experiment #1 (v0.1.0)
_August 19th, 2024_

## Introduction

The first iteration should be considered a proof of concept. There was not a thoughtful process behind the aesthetics of the produced sounds but rather the focus was to come up with a working strategy of generating them. Further experiments will address the shortcomings of the first one.

Experiment #1 comes in the form of a small command line program that generates audio based on markov chains models. It has five available flags. `-debug` to enable debug logging - useful while developing, `-ngen` the number of successive generations to generate, `-files` the directory to save the audio files of each generation, `-models` the directory to save the generated markov chains models and `-seed` the seed markov chain models to kick start the process.

Running it produces folders containing the generated audio and the Markov chain models used for the generation in json format.

```
├── gen0
│   ├── 0.wav
│   ├── 1.wav
│   ├── amp.json
│   ├── dur.json
│   ├── freq.json
├── gen1
│   ├── 0.wav
│   ├── 1.wav
│   ├── amp.json
│   ├── dur.json
│   ├── freq.json
├── gen2
│   ├── 0.wav
│   ├── 1.wav
│   ├── amp.json
│   ├── dur.json
│   ├── freq.json
...
```

_Note that to run the program you need to either install `go` or download the binary named `markov1` from the module's github under [releases](https://github.com/bh90210/mlsic/releases)._

Under the hood the CLI uses the `markov` [package](https://github.com/bh90210/mlsic/tree/trunk/markov) to generate a "train" of sines waves. Experiment #1 uses additive synthesis to produce audio signal. The result is a monophonic synth. Each fundamental is treated for harmonics. Harmonics are read of a corresponding table and for iteration #1 are static and the same for each fundamental. Each generation uses the previous generation's models to generate new values for the sine waves. Gen0 uses the seed models.

{{< mermaid class="optional" >}}
flowchart TD
    a["Ngen"] -->
    c["freq.json"] --> 
    d{"`**Markov Generator**`"} -->
    g[Create Train]
    a -- Load the Markov models ---
    e["amp.json"] -->
    d
    g
    
    a --> 
    h["dur.json"] -->
    d
    g

    g -->
    k["Harmonics"]

    k -->
    l["`**Generate Audio**`"]

    l -->
    p["Generate New Models"]

    p -.-> a;
{{< /mermaid >}}

## Strategy

Experiment #1 uses Markov Chains to generate variations of the initial seed. The starting point is the three seed models (frequencies, for amplitudes, durations.) Each model is read and fed into a Markov generator. 

```json
		"0.000000": 137,
		"1.000000": 206,
		"10.000000": 94,
		"10.148148": 123,
		"10.232558": 93,
		"10.301887": 122,
		"10.461538": 121,
		"10.476190": 92,
		"10.627451": 120,
		"10.731707": 91,
		"10.800000": 119,
```
_Sample values from the seed freq.json model._

{{< mermaid class="optional" >}}
flowchart TD
    a["`**Freq Model**
    	0.000000
		1.000000
		10.000000
		10.148148
		10.232558
		10.301887
		10.461538
		10.476190
		10.627451
		10.731707
		10.800000
        ...`"]
    b{Markov Generator}

    a -- Load Model --- b
{{< /mermaid >}}

Each individual value of the model is fed into the generator until either a. the generator stops producing new values or b. the produced values stop being unique and start looping (e.g. `0.000000, 1.000000, 10.000000, 0.000000, 1.000000, 10.000000 ...`.) 

{{< mermaid class="optional" >}}
flowchart TD
    a{Markov Generator}

    b["`**Freq Model**
    	0.000000
        ...`"]

    b -- Feed First Value --- a

    a --> c["`**Generated Values**
    10.148148
    2100.000000
    19.130435
    20.952381
    2750.000000
    ...`"]
{{< /mermaid >}}

This produces an array of float values. Once the subprocess is over the generator is fed the next value of the model `1.000000` and generates a new array of float values. The process continues until we feed all values of the model to the generator.

The end result is three arrays of arrays `[][]float`. Then the program reads through all generated values of frequencies, amplitudes and durations and creates a train of Sines.

```golang
type Sine struct {
	Frequency float64
	Amplitude float64
	Duration  time.Duration
}
```

{{< mermaid class="optional" >}}
flowchart TD
    a["`Freqs
     [0][0.000000, 1.000000 ...]
     [1][10.76190, 19.13043 ...]
     ...`"]
    b["`Amps 
     [0][0.510000, 0.350000 ...]
     [1][0.450000, 0.510000 ...]
     ...`"]
    c["`Durs 
     [0][129, 130 ...]
     [1][143, 150 ...]
     ...`"]

    d{Sine Constructor}

    a & b & c --> d


    e["`Sine Train
    Sine {
    Frequency: Freqs[0][0]
	Amplitude: Amps[0][0]
	Duration: Durs[0][0]}
    Sine {
	Frequency: Freqs[0][1]
	Amplitude: Amps[0][1]
	Duration: Durs[0][1]}
    Sine {
	Frequency: Freqs[0][2]
	Amplitude: Amps[0][2]
	Duration: Durs[0][2]}
...`"]

    d --> e
{{< /mermaid >}}

At this point we have a representation of monophonic consecutive pure sine waves. The next step is to generate the audio signal based on them. Along the fundamentals we generate the partials for each and add them together.

```
2 0.02
3 0.03
4 0.04
5 0.05
6 0.06
7 0.07
8 0.08
9 0.09
10 0.1
11 0.11
12 0.12
...
```
_Excerpt of the harmonics table. First column is the partial (the 2nd, the 3rd etc) and the second column is the factor to multiply against the fundamental amplitude to derive the amplitude of the partial._

The final step is to save the generated audio files and export the new models that will be used as seeds for the next n generation.

## Seed

<audio src="https://bh90210.github.io/mlsic/blog/experiment_1_seed.wav" controls preload></audio>

## Generator

## Result

<audio src="/mlsic/blog/experiment_1_ngen_24_left_channel_ngen_16 both_channels_ngen_40_right_channel_excerpt.wav" controls preload></audio>

<audio src="/mlsic/blog/experiment_1_ngen_25_left_channel_ngen_39_right_channel.wav" controls preload></audio>

#### Limitations

The code so far is very crude and suffers from many limitations. Limitations here are understood as impediments to reaching a richer and/or more controlled audio signal. The list bellow is not exhaustive but will be used as a guide for Experiment #2.

* Crude seed

The way the algorithm works the initial seed has the biggest influence on the end result regardless of how many generations past it we are.
The two solutions are a. more complex seeds to begin with, b. a new strategy for combining multiple seeds in the generation process.

* Monophonic

The algorithm produces mono signal (one channel.) We managed to achieve stereo image by combining three generations and panning them. Modification are needed to auto generate stereo, quadraphonic audio etc.

* Monophony

The procedure of the audio signal generation is such that the end result can not be polyphonic. Modification are needed to achieve polyphony.

* Static harmonics

The same harmonics table is applied to all Sines of the train. A proper strategy is needed for dynamic harmonics.

* Generation sequence

Each model's values are sorted (see _Sample values from the seed freq.json model._ above.) For example value `0.000000` is fed to the Markov Generator. When the process finishes value `1.000000` is fed and so on. This has profound effect on the audio signal generated. A new strategy is needed to make this process more dynamic.

* Substitutions in Sine Constructor

As described previously, the Markov Generator is used to produce three arrays of arrays. Those are fed to the Sine Constructor. At the moment Constructor favours heavily the generated result of frequencies. This happens in two major ways. When a frequencies sub array (eg. `Freqs[0]`) reaches the end, it ignores the remaining amplitudes and durations arrays (if any.) While reading the `Freqs[0]` sub array if `len(Freqs[0] > len(Amps[0])`, that is there is no corresponding amplitude value for the next frequency of the Sine Train, hardcoded default value 0 is used (no signal.) Same applies for duration, default hard coded value when there are no duration value left is 10 milliseconds.
