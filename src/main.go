package main

import (
	"image"

	"./palette"
	"./parallel"
)

func main() {

	colorPalette := palette.Vivid()
	imageFrame := image.Rect(0, 0, 1280, 1024)
	parallel.Run(colorPalette, 500, imageFrame, -0.6366988, -0.4426395)
}
