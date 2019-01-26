package main

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"math/cmplx"
	"os"
	"strconv"

	"../palette"
)

const (
	screenWidth    = 1280
	aspectRatio    = 5 / 4
	noOfIterations = 1000
	escapeRadius   = 200
)

func main() {

	realMin := -2.
	realMax := .5
	imagMin := -1.
	imagMax := 1.
	counter := 0
	scale := screenWidth / (realMax - realMin)
	screenHeight := int(scale * (imagMax - imagMin))
	colorPalette := palette.Vivid()

	for counter < 1 {
		bounds := image.Rect(0, 0, screenWidth, screenHeight)
		img := image.NewPaletted(bounds, colorPalette)
		fmt.Println(scale)
		for x := 0; x < screenWidth; x++ {
			for y := 0; y < screenHeight; y++ {
				i := mandelbrot(complex(float64(x)/scale+realMin, float64(y)/scale+imagMin))
				colIndex := uint8(i * float64(len(colorPalette)-1))
				img.SetColorIndex(x, y, colIndex)
			}
		}
		f, _ := os.Create("mandelbrot" + strconv.Itoa(counter) + ".png")
		png.Encode(f, img)
		f.Close()
		counter++
		scale *= 1.05

	}
}

func mandelbrot(c complex128) float64 {
	iteration := 0
	z := c
	for cmplx.Abs(z) < escapeRadius && iteration < noOfIterations {
		z = z*z + c
		iteration++
	}
	if iteration == noOfIterations {
		return 1.0
	}
	return (float64(iteration) + 1.0 - math.Log(math.Log2(float64(cmplx.Abs(z))))) / noOfIterations
}
