package sequential

import (
	"image/color"
)

import (
	"fmt"
	"image"
	"image/gif"
	"math"
	"math/cmplx"
	"os"

	"../palette"
)

const (
	screenWidth    = 1280
	aspectRatio    = 5 / 4
	noOfIterations = 1000
	escapeRadius   = 200
)

type work struct {
	scale
	realMin
	imagMin
}

func Run(colorPalette color.Palette, frame int) {

	finishedTasks := func(ch chan<- work) chan work {
		done := make(chan work, frames)
		for i:= 0; i < runtime.RunCPU(); i++{
		go workThread(done)
		}
	}
}

	outGIF := &gif.GIF{}

	realMin := -2.
	realMax := .5
	imagMin := -1.
	imagMax := 1.
	counter := 0
	scale := screenWidth / (realMax - realMin)
	screenHeight := int(scale * (imagMax - imagMin))
	colorPalette := palette.Vivid()

	for counter < 40 {
		bounds := image.Rect(0, 0, screenWidth, screenHeight)
		img := image.NewPaletted(bounds, colorPalette)
		for x := 0; x < screenWidth; x++ {
			for y := 0; y < screenHeight; y++ {
				i := mandelbrot(complex(float64(x)/scale+realMin, float64(y)/scale+imagMin))
				colIndex := uint8(i * float64(len(colorPalette)-1))
				img.SetColorIndex(x, y, colIndex)
			}
		}

		outGIF.Image = append(outGIF.Image, img)
		outGIF.Delay = append(outGIF.Delay, 0)
		counter++

		realMin += (1. / (scale * 0.025))
		imagMin += (1. / (scale * 0.025))
		scale *= 1.05
	}

	f, _ := os.Create("giffy.gif")
	gif.EncodeAll(f, outGIF)
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
