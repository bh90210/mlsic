<p align="left">
  <a href="https://goreportcard.com/report/github.com/bh90210/mlsic">
    <img src="https://goreportcard.com/badge/github.com/bh90210/mlsic" />
  </a>
  <a href="https://pkg.go.dev/github.com/bh90210/mlsic">
    <img src="https://pkg.go.dev/badge/github.com/bh90210/mlsic.svg" alt="Go Reference"/>
  </a>
  <a href="https://codecov.io/gh/bh90210/mlsic" > 
 <img src="https://codecov.io/gh/bh90210/mlsic/branch/trunk/graph/badge.svg?token=ZT1OQEETCQ"/> 
 </a>
  <a href="https://github.com/bh90210/mlsic/actions/workflows/ci.yml" rel="nofollow">
    <img src="https://img.shields.io/github/actions/workflow/status/bh90210/mlsic/ci.yml?branch=trunk&logo=Github" alt="Build" />
  </a>
    <a href="https://github.com/bh90210/mlsic/blob/trunk/LICENSE">
    <img alt="GitHub" src="https://img.shields.io/github/license/bh90210/mlsic"/>
  </a>
</p>


# Mlsic v1

IN DEVELOPMENT

# How to get started

First take a look in the `/example/` directory for a quick high use example of the module.

#### Go

You need [Go 1.19 installed](https://go.dev/doc/install) on your system.

#### PortAudio

You also need [PortAudio](http://www.portaudio.com/) installed.

On Debian based linux distors you can do `apt install portaudio19-dev`.

On MacOS you can use Hoembrew `brew install portaudio`.

Alternatively under `/render` you will find a wav and aiff renderer for 100% offline audio generation.

#### Dot cli

If you wish to render SVG files out of your generated graphs (see `/generator`) you will need Graphviz Dot [cli installed](https://graphviz.org/download/). 