package parallel

import (
	"image"
	"image/color"
	"image/gif"
	"math"
	"math/cmplx"
	"os"
	"runtime"
	"sync"
)

const (
	noOfIterations = 1000
	escapeRadius   = 200
)

type work struct {
	frameNo int
	scale   float64
	realMin float64
	imagMin float64
	wg      *sync.WaitGroup
}

func Run(colorPalette color.Palette, frames int, imageFrame image.Rectangle) {

	outGIF := &gif.GIF{}
	images := make([]*image.Paletted, frames)

	workPile := make(chan work, frames)
	for i := 0; i < runtime.NumCPU(); i++ {
		go func(workChannel <-chan work) {
			for {
				work, hasMore := <-workChannel
				if hasMore {
					img := image.NewPaletted(imageFrame, colorPalette)
					for x := 0; x < imageFrame.Max.X; x++ {
						for y := 0; y < imageFrame.Max.Y; y++ {

							i := mandelbrot(complex(float64(x)/work.scale+work.realMin, float64(y)/work.scale+work.imagMin))
							colIndex := uint8(i * float64(len(colorPalette)-1))
							img.SetColorIndex(x, y, colIndex)
						}
					}
					images[work.frameNo] = img
					work.wg.Done()
				} else {
					break
				}
			}
		}(workPile)
	}

	realMin := -2.
	realMax := .5
	imagMin := -1.
	scale := float64(imageFrame.Max.X) / (realMax - realMin)

	var wg sync.WaitGroup
	for i := 0; i < frames; i++ {
		wg.Add(1)
		work := work{frameNo: i, scale: scale, realMin: realMin, imagMin: imagMin, wg: &wg}
		workPile <- work
		realMin += (1. / (scale * 0.07))
		imagMin += (1. / (scale * 0.07))
		scale *= 1.02
	}

	outGIF.Delay = make([]int, frames)

outer:
	for {
		select {
		case work := <-workPile:
			img := image.NewPaletted(imageFrame, colorPalette)
			for x := 0; x < imageFrame.Max.X; x++ {
				for y := 0; y < imageFrame.Max.Y; y++ {

					i := mandelbrot(complex(float64(x)/work.scale+work.realMin, float64(y)/work.scale+work.imagMin))
					colIndex := uint8(i * float64(len(colorPalette)-1))
					img.SetColorIndex(x, y, colIndex)
				}
			}
			images[work.frameNo] = img
			work.wg.Done()
		default:
			wg.Wait()
			break outer
		}
	}

	wg.Wait()

	outGIF.Image = images
	file, _ := os.Create("poop.gif")
	gif.EncodeAll(file, outGIF)
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
