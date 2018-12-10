package main

import (
	"fmt"
	"image"

	"./opencl"
	"./palette"
)

func main() {

	colorPalette := palette.Vivid()
	fmt.Println(len(colorPalette))
	imageFrame := image.Rect(0, 0, 1440, 1152)
	//parallel.Run(colorPalette, 500, imageFrame, -0.6366988, -0.4426395)
	opencl.Run(colorPalette, 10, imageFrame, -0.6366988, -0.4426395)
}
