package main

import (
	"./Vectores"
	"image/color"
)

var FL float64 = 0.5

const CTEOCTAEDRO = 0.57735027

var EYE = Vectores.Vector{0, 0, 0}
var BACKGROUNDCOLOR = color.RGBA{R: 0, G: 0, B: 0, A: 255}
var NOISECOLOR = color.RGBA{R: 250, G: 250, B: 250, A: 255}

var WIDTH int = 800
var HEIGHT int = 600

// Ángulo para el FOV. Actúa como una especie de zoom.
var ALPHA float32 = 55.0

var correccion float64 = 0.5
var ImageAspectRatio float64 = float64(WIDTH) / float64(HEIGHT)
var MAXSTEPS = 32
var MINIMUM_HIT_DISTANCE = 0.35

var currentColor color.RGBA

var CurrentMaterial int

const STEP float64 = .1
const RADIO_ESFERA float64 = 5
const NOISEZOOM = 0.25
