package opencl

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"log"
	"math"
	"os"
	"unsafe"

	"github.com/jgillich/go-opencl/cl"
)

var kernelSource = `
__kernel void mandelbrot(__global char *out, const float realMin, const float realMax, const float imagMin, const float imagMax, const int paletteCount) {
	int x_dim = get_global_id(0);
	int y_dim = get_global_id(1);
	size_t width = get_global_size(0);
	size_t height = get_global_size(1);
	
	int pixel = y_dim * width + x_dim;

	float c_x = ((float) x_dim / width)*(realMax - realMin) + realMin;
	float c_y = ((float) y_dim / height)*(imagMax - imagMin) + imagMin;

	float x = 0.0;
	float y = 0.0;

	int iteration = 0;
	int noOfIterations = 1000;

	while(x * x + y * y <= 200 && iteration < noOfIterations) {
		float xtemp = x * x - y * y + c_x;
		y = 2 * x * y + c_y;
		x = xtemp;
		iteration++;
	}

  	if(iteration == noOfIterations) {
		out[pixel] = paletteCount - 1;
  	} else {
	float fraction = ((float) iteration + 1.0 - log(log2((x * x + y * y)))) / noOfIterations;
	out[pixel] = fraction * (paletteCount - 1);
	}
}
`

func Run(colorPalette color.Palette, frames int, imageFrame image.Rectangle, destX float64, destY float64) {
	imageWidth := imageFrame.Max.X
	imageHeight := imageFrame.Max.Y

	platforms, err := cl.GetPlatforms()
	if err != nil {
		log.Fatal(err)
	}
	platform := platforms[0]
	devices, err := platform.GetDevices(cl.DeviceTypeGPU)
	if err != nil {
		log.Fatal(err)
	}
	if len(devices) == 0 {
		log.Fatalf("Could not find any GPUs")
	}

	availableDevices := make([]*cl.Device, 0)
	for _, device := range devices {
		if device.Available() {
			fmt.Println(device.Name())
			availableDevices = append(availableDevices, device)
		}
	}

	context, err := cl.CreateContext(availableDevices)
	if err != nil {
		log.Fatal(err)
	}

	cmdQueues := make([]*cl.CommandQueue, len(availableDevices))

	for i, device := range devices {
		cmdQueues[i], err = context.CreateCommandQueue(device, 0)
		if err != nil {
			log.Fatal(err)
		}
	}

	program, err := context.CreateProgramWithSource([]string{kernelSource})
	if err != nil {
		log.Fatal(err)
	}

	if err := program.BuildProgram(nil, ""); err != nil {
		log.Fatal(err)
	}

	kernel, err := program.CreateKernel("mandelbrot")
	if err != nil {
		log.Fatal(err)
	}

	realMin := float32(-2.)
	realMax := float32(.5)
	imagMin := float32(-1.)
	imagMax := float32(1.)
	scale := float32(imageWidth) / (realMax - realMin)
	easeOut := float32(0.98)
	dX := float32(math.Abs(destX-float64(realMin)) * (1. - float64(easeOut)) / (1. - math.Pow(float64(easeOut), float64(frames))))
	dY := float32(math.Abs(destY-float64(imagMin)) * (1. - float64(easeOut)) / (1. - math.Pow(float64(easeOut), float64(frames))))
	images := make([]*image.Paletted, frames)
	delay := make([]int, frames)
	outGIF := &gif.GIF{Image: images, Delay: delay}

	for i := 0; i < frames; i++ {

		outBuffer, err := context.CreateEmptyBuffer(cl.MemWriteOnly, imageHeight*imageWidth)
		if err != nil {
			log.Fatal(err)
		}

		err = kernel.SetArgs(outBuffer, realMin, realMax, imagMin, imagMax, int32(len(colorPalette)))
		if err != nil {
			log.Fatal(err)
		}

		img := image.NewPaletted(imageFrame, colorPalette)
		images[i] = img
		workPerQueue := []int{imageWidth, imageHeight / len(cmdQueues)}
		for i := 0; i < len(cmdQueues); i++ {
			queueOffset := []int{0, workPerQueue[1] * i}
			absoluteOffset := queueOffset[1] * imageWidth

			_, err := cmdQueues[i].EnqueueNDRangeKernel(kernel, queueOffset, []int{1440, 1152}, nil, nil)
			if err != nil {
				log.Fatal(err)
			}
			_, err = cmdQueues[i].EnqueueReadBuffer(outBuffer, false, absoluteOffset, workPerQueue[0]*workPerQueue[1], unsafe.Pointer(&img.Pix[absoluteOffset]), nil)
			if err != nil {
				log.Fatal(err)
			}
		}

		realMin += dX
		imagMin += dY
		scale *= 1.02
		realMax = realMin + float32(imageWidth)/scale
		imagMax = imagMin + float32(imageHeight)/scale
		dX *= easeOut
		dY *= easeOut
	}

	for i := 0; i < len(cmdQueues); i++ {
		cmdQueues[i].Finish()
	}

	f, _ := os.Create("giffy.gif")
	gif.EncodeAll(f, outGIF)
}
