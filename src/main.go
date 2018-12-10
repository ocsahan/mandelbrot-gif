package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"

	"./opencl"
	"./palette"
	"./parallel"
	"./sequential"
)

func main() {

	var colors = flag.String("colors", "Vivid", "Refer to README for color options. Default is Vivid.")
	var frame = flag.Int("frame", 100, "Number of frames to get from (0,0) to dest. Default is 100.")
	var resolution = flag.String("resolution", "medium", "Image size. Values: [small, medium, large]. Default is medium.")
	var destX = flag.Float64("dest_x", -0.6366988, "Origin of the last frame in coordinate plane. Default is -0.6366988")
	var destY = flag.Float64("dest_y", -0.4426395, "Origin of the last frame in coordinate plane. Default is -0.4426395")
	var parallelism = flag.Bool("p", false, "Parallelism flag. True or False. Default is False.")
	var threads = flag.Int("threads", -1, "Number of threads. Can only be used with -p=true flag. Value 0 will use the GPU.")

	flag.Parse()
	flag.PrintDefaults()

	var colorPalette color.Palette
	switch *colors {
	case "Vivid":
		colorPalette = palette.Vivid()
	case "Hippie":
		colorPalette = palette.Hippie()
	}

	var imageFrame image.Rectangle
	switch *resolution {
	case "large":
		imageFrame = image.Rect(0, 0, 1440, 1152)
	case "medium":
		imageFrame = image.Rect(0, 0, 1040, 832)
	case "small":
		imageFrame = image.Rect(0, 0, 640, 512)
	}

	if *parallelism {
		switch *threads {
		case -1:
			parallel.Run(colorPalette, *frame, imageFrame, *destX, *destY, *threads)
		case 0:
			opencl.Run(colorPalette, *frame, imageFrame, *destX, *destY)
			fmt.Println(*parallelism)
			fmt.Println(*threads)
		default:
			parallel.Run(colorPalette, *frame, imageFrame, *destX, *destY, *threads)
		}
	} else {
		if *threads != -1 {
			flag.PrintDefaults()
			return
		}
		sequential.Run(colorPalette, *frame, imageFrame, *destX, *destY)
	}
}
