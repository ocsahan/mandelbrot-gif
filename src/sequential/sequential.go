package sequential

import (
	"image"
	"image/color"
	"image/gif"
	"math"
	"math/cmplx"
	"os"
)

const (
	screenWidth    = 1280
	aspectRatio    = 5 / 4
	noOfIterations = 1000
	escapeRadius   = 200
)

func Run(colorPalette color.Palette, frames int, imageFrame image.Rectangle, destX float64, destY float64) {

	outGIF := &gif.GIF{}
	screenWidth := imageFrame.Max.X
	screenHeight := imageFrame.Max.Y

	realMin := -2.
	realMax := .5
	imagMin := -1.
	scale := float64(imageFrame.Max.X) / (realMax - realMin)
	easeOut := 0.98
	dX := math.Abs(destX-realMin) * (1. - easeOut) / (1. - math.Pow(easeOut, float64(frames)))
	dY := math.Abs(destY-imagMin) * (1. - easeOut) / (1. - math.Pow(easeOut, float64(frames)))

	for frame := 0; frame < frames; frame++ {
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

		realMin += dX
		imagMin += dY
		dX *= easeOut
		dY *= easeOut
		scale *= 1.02
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
