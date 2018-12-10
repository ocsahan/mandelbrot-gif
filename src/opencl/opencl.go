package opencl

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"unsafe"

	"../utils"

	"github.com/jgillich/go-opencl/cl"
)

var kernelSource = `
__kernel void mandelbrot(__global char *out, const float realMin, const float realMax, const float imagMin, const float imagMax, const int paletteCount) {
	int x_dim = get_global_id(0);
	int y_dim = get_global_id(1);
	size_t width = get_global_size(0);
	size_t height = get_global_size(1);
	int idx = y_dim * width + x_dim;

	float x_origin = ((float) x_dim / width)*(realMax - realMin) + realMin;
	float y_origin = ((float) y_dim / height)*(imagMax - imagMin) + imagMin;

  // The Escape time algorithm, it follows the pseduocode from Wikipedia
  // _very_ closely
  float x = 0.0;
  float y = 0.0;

  int iteration = 0;
  int noOfIterations = 1000;

  while(x*x + y*y <= 200 && iteration < noOfIterations) {
    float xtemp = x*x - y*y + x_origin;
    y = 2*x*y + y_origin;
    x = xtemp;
    iteration++;
  }

  if(iteration == noOfIterations) {
    // This coordinate did not escape, so it is in the Mandelbrot set
    out[idx] = paletteCount - 1;
  } else {
    // This coordinate did escape, so color based on quickly it escaped
	float fraction = ((float) iteration + 1.0 - log(log2((x*x + y*y)))) / noOfIterations;
	out[idx] = fraction * (paletteCount-1);
  }
}
`

func Run(colorPalette color.Palette, frames int, imageFrame image.Rectangle, coords utils.CoordFrame, destX float64, destY float64) {
	imageWidth := imageFrame.Max.X
	imageHeight := imageFrame.Max.Y

	platforms, err := cl.GetPlatforms()
	if err != nil {
		log.Fatal(err)
	}
	platform := platforms[0]
	devices, err := platform.GetDevices(cl.DeviceTypeAll)
	if err != nil {
		log.Fatal(err)
	}
	if len(devices) == 0 {
		log.Fatalf("Could not find any devices")
	}

	context, err := cl.CreateContext(devices)
	if err != nil {
		log.Fatal(err)
	}

	cmdQueues := make([]*cl.CommandQueue, len(devices))

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

	outBuffer, err := context.CreateEmptyBuffer(cl.MemWriteOnly, imageHeight*imageWidth)
	if err != nil {
		log.Fatal(err)
	}

	err = kernel.SetArgs(outBuffer, coords.RealMin, coords.RealMax, coords.ImagMin, coords.ImagMax, int32(len(colorPalette)))
	if err != nil {
		log.Fatal(err)
	}

	img := image.NewPaletted(imageFrame, colorPalette)
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

	for i := 0; i < len(cmdQueues); i++ {
		cmdQueues[i].Finish()
	}

	f, err := os.Create("image.png")
	png.Encode(f, img)
}
