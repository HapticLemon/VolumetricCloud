package main

import (
	"./Ruido"
	"./Vectores"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"
	"time"
)

// Para simplificar, considero una esfera de radio1 en el centro de coordenadas.
//
func distanciaEsfera(punto Vectores.Vector) float64 {
	// 10 es el radio de la esfera.
	//var ruido float64 = fbm(punto, 1)

	//var translation = Vectores.Vector{0 + ruido, 0 + ruido, -22.0}
	//var translation = Vectores.Vector{0 , 0 , -22.0}
	var translation = Ruido.CurlNoise(punto.MultiplyByScalar(0.5))
	translation = Vectores.Vector{translation.X * 3, translation.Y * 3, -20}

	return punto.Sub(translation).Length() - RADIO_ESFERA

}

func sumaOctavas(num_iterations int, punto Vectores.Vector, persistence float64, scale float64, low float64, high float64) float64 {
	var maxamp float64 = 0
	var amp float64 = 1
	var freq float64 = scale
	var noise float64 = 0

	for i := 0; i < num_iterations; i++ {
		noise += Ruido.Noise3(punto.X*NOISEZOOM, punto.Y*NOISEZOOM, punto.Z*NOISEZOOM) * amp
		maxamp += amp
		amp *= persistence
		freq *= 2
	}

	noise /= maxamp
	noise = noise*(high-low)/2 + (high+low)/2
	return noise
}

func calculaDensidadLineal(punto Vectores.Vector, longitud float64) float64 {
	var noiseValue float64 = 0

	//punto.MultiplyByScalar(NOISEZOOM)

	noiseValue = (Ruido.Noise3(punto.X*NOISEZOOM, punto.Y*NOISEZOOM, punto.Z*NOISEZOOM)) * 0.015
	//noiseValue = Ruido.Worley3D(punto) * 0.0025
	//noiseValue = sumaOctavas(16, punto, .5, 0.007, 0, 255) * 0.00005
	return math.Abs(noiseValue * longitud)
}

// Implementación de niebla según idea de Íñigo Quílez.
// https://iquilezles.org/www/articles/fog/fog.htm
func applyFog(color color.RGBA, distancia float64, densidad float64) color.RGBA {
	var fogAmount float32 = 0.0

	fogAmount = float32(1.0 - math.Pow(math.E, -distancia*densidad)) //* 1.5

	return mixColor(color, BACKGROUNDCOLOR, fogAmount)
	//return mixColor(BACKGROUNDCOLOR, color, fogAmount)
}

// Interpolación entre dos colores.
//
func mixColor(x color.RGBA, y color.RGBA, a float32) color.RGBA {
	var resultado color.RGBA

	resultado.R = uint8(float32(x.R)*(1-a) + float32(y.R)*a)
	resultado.G = uint8(float32(x.G)*(1-a) + float32(y.G)*a)
	resultado.B = uint8(float32(x.B)*(1-a) + float32(y.B)*a)

	return resultado
}

// Implementación de FBM basada en
// https://www.iquilezles.org/www/articles/fbm/fbm.htm
//
func fbm(punto Vectores.Vector, h float64) float64 {
	var G float64 = math.Exp2(-h)
	var f float64 = 1.0
	var a float64 = 1
	var t float64 = 0

	for i := 0; i < OCTAVES; i++ {
		t += a * Ruido.Noise3(f*punto.X, f*punto.Y, f*punto.Z)
		f *= 2
		a *= G
	}
	return t
}

func raymarch(ro Vectores.Vector, rd Vectores.Vector) color.RGBA {

	var punto Vectores.Vector
	//var puntoCurl Vectores.Vector
	var t float64 = 0
	var densidad float64 = 0
	var longitud float64 = 0

	var color = BACKGROUNDCOLOR

	for x := 0; x < MAXSTEPS; x++ {
		punto = ro.Add(rd.MultiplyByScalar(t))
		//punto = Ruido.CurlNoise(punto)
		distancia := distanciaEsfera(punto)

		// Hemos tocado la esfera.
		// Aplico ruído a la distancia mínima para distorsionar el contorno de la esfera.
		//if distancia < MINIMUM_HIT_DISTANCE-fbm(punto, 1.5) {
		if distancia < MINIMUM_HIT_DISTANCE {
			//return NOISECOLOR
			//puntoCurl = Ruido.CurlNoise(punto)
			//distancia = distanciaEsfera(puntoCurl)
			for distancia < RADIO_ESFERA {
				//densidad += calculaDensidadLineal(punto, longitud)
				densidad += calculaDensidadLineal(punto, STEP)
				longitud += STEP
				punto = ro.Add(rd.MultiplyByScalar(t + longitud))
				distancia = distanciaEsfera(punto)

			}
			return applyFog(NOISECOLOR, longitud, densidad)
		}
		t += distancia
	}

	// Devuelvo el color de fondo.
	return color
}

func main() {
	var NDC_x float64
	var NDC_y float64
	var PixelScreen_x float64
	var PixelScreen_y float64
	var PixelCamera_x float64
	var PixelCamera_y float64

	var ro Vectores.Vector
	var rd Vectores.Vector
	var nuevo Vectores.Vector
	var color color.RGBA

	var fileOut string

	start := time.Now()

	argsWithoutProg := os.Args[1:]

	fileOut = argsWithoutProg[0]

	fmt.Printf("Files Out %s\n", fileOut)

	img := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
	out, err := os.Create(fileOut)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Calculo el Field of View. El ángulo es de 45 grados.
	//
	var FOV float64 = float64(math.Tan(float64(ALPHA / 2.0 * math.Pi / 180.0)))

	for x := 0; x < WIDTH; x++ {
		for y := 0; y < HEIGHT; y++ {
			// Hacemos las conversiones de espacios
			//
			NDC_x = (float64(x) + correccion) / float64(WIDTH)
			NDC_y = (float64(y) + correccion) / float64(HEIGHT)

			PixelScreen_x = 2*NDC_x - 1
			PixelScreen_y = 2*NDC_y - 1

			PixelCamera_x = PixelScreen_x * ImageAspectRatio * FOV
			PixelCamera_y = PixelScreen_y * FOV

			// Origen y dirección

			ro = EYE
			nuevo.X = PixelCamera_x
			nuevo.Y = PixelCamera_y
			nuevo.Z = -1

			rd = nuevo.Sub(ro).Normalize()

			color = raymarch(ro, rd)

			img.Set(x, y, color)

		}
	}
	elapsed := time.Since(start)
	log.Printf("Binomial took %s", elapsed)

	var opt jpeg.Options

	opt.Quality = 80
	// ok, write out the data into the new JPEG file

	err = jpeg.Encode(out, img, &opt) // put quality to 80%
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Generated image to %s \n", fileOut)

}
