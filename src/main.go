package main

import (
	"fmt"
	"image"

	"./opencl"
	"./utils"
)

func main() {

	colorPalette := utils.Vivid()
	fmt.Println(len(colorPalette))
	imageFrame := image.Rect(0, 0, 1440, 1152)
	//parallel.Run(colorPalette, 500, imageFrame, -0.6366988, -0.4426395)
	coords := utils.CoordFrame{RealMin: -2., RealMax: .5, ImagMin: -1, ImagMax: 1}
	opencl.Run(colorPalette, 1, imageFrame, coords, 0, 0)
}
