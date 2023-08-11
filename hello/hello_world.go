package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	pb "github.com/cheggaaa/pb/v3"
	glm "github.com/engoengine/glm"
)

func main() {
	const width, height = 640, 480

	fmt.Println("Go Raytracer!")

	image := image.NewNRGBA(image.Rect(0, 0, width, height))

	bar := pb.StartNew(height)

	for y := 0; y < height; y++ {
		bar.Increment()
		for x := 0; x < width; x++ {

			pixel := glm.Vec3{(float32(x) / float32(width-1)), (float32(y) / float32(height-1)), float32(0)}

			writeColours(image, x, y, pixel)
		}
	}
	bar.Finish()

	f, err := os.Create("image.png")
	if err != nil {
		fmt.Println(err)
	}

	if err := png.Encode(f, image); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		fmt.Println(err)
	}
}

func writeColours(image *image.NRGBA, x, y int, pixel glm.Vec3) {
	ir := uint8(255.999 * pixel.X())
	ig := uint8(255.999 * pixel.Y())
	ib := uint8(255.999 * pixel.Z())

	image.Set(x, y, color.NRGBA{
		R: ir,
		G: ig,
		B: ib,
		A: 255,
	})
}
