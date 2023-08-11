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

type at func(float32) glm.Vec3
type origin func() glm.Vec3
type Direction func() glm.Vec3

type Ray struct {
	Origin    glm.Vec3
	Direction glm.Vec3

	At   at
	Orig origin
	Dir  Direction
}

func newRay(o, d glm.Vec3) *Ray {
	r := new(Ray)

	r.Origin = o
	r.Direction = d

	r.At = func(t float32) glm.Vec3 {
		temp := r.Direction.Mul(t)
		or := r.Orig()
		return or.Add(&temp)
	}

	r.Orig = func() glm.Vec3 {
		return r.Origin
	}

	r.Dir = func() glm.Vec3 {
		return r.Direction
	}
	return r
}

func main() {
	const aspectRatio = 16.0 / 9.0
	const image_width = 1920
	const image_height = int(image_width / aspectRatio)

	focal_length := 1.0

	camera_center := glm.Vec3{0.0, 0.0, 0.0}

	const viewportHeight = 2.0
	const viewportWidth = (float32(image_width) / float32(image_height)) * viewportHeight

	viewport_u := glm.Vec3{viewportWidth, 0, 0}
	viewport_u_half := viewport_u.Mul(0.5)

	viewport_v := glm.Vec3{0, -viewportHeight, 0}
	viewport_v_half := viewport_v.Mul(0.5)

	pixel_delta_u := viewport_u.Mul(1.0 / float32(image_width))
	pixel_delta_v := viewport_v.Mul(1.0 / float32(image_height))

	focal_vec := glm.Vec3{0, 0, float32(focal_length)}

	viewport_upper_left_corner := camera_center.Sub(&focal_vec)
	viewport_upper_left_corner = viewport_upper_left_corner.Sub(&viewport_u_half)
	viewport_upper_left_corner = viewport_upper_left_corner.Sub(&viewport_v_half)

	pixel_sum := pixel_delta_u.Add(&pixel_delta_v)
	pixel_sum = pixel_sum.Mul(0.5)
	pixel00_loc := viewport_upper_left_corner.Add(&pixel_sum)

	fmt.Println("Go Raytracer!")

	image := image.NewNRGBA(image.Rect(0, 0, image_width, image_height))

	bar := pb.StartNew(image_height)

	for y := 0; y < image_height; y++ {
		bar.Increment()
		for x := 0; x < image_width; x++ {
			pixel_u_x := pixel_delta_u.Mul(float32(x))
			pixel_v_y := pixel_delta_v.Mul(float32(y))
			pixel_center := pixel00_loc.Add(&pixel_u_x)
			pixel_center = pixel_center.Add(&pixel_v_y)

			ray_dir := pixel_center.Sub(&camera_center)

			r := newRay(camera_center, ray_dir)

			pixel := RayColour(r)

			writeColours(image, x, y, pixel)
		}
	}
	bar.Finish()

	fmt.Println("Finish")

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

func RayColour(r *Ray) glm.Vec3 {

	if hit_sphere(glm.Vec3{0, 0, -1}, 0.5, r) {
		return glm.Vec3{1, 0, 0}
	}

	ray_dir := r.Dir()
	unit_dir := ray_dir.Normalized()

	a := 0.5*unit_dir.Y() + 1.0
	white := glm.Vec3{1.0, 1.0, 1.0}
	white = white.Mul(1.0 - a)

	blue := glm.Vec3{0.5, 0.7, 1.0}
	blue = blue.Mul(a)

	return white.Add(&blue)
}

func hit_sphere(center glm.Vec3, radius float32, r *Ray) bool {
	oc := r.Origin.Sub(&center)
	a := r.Direction.Dot(&r.Direction)
	b := 2.0 * oc.Dot(&r.Direction)
	c := oc.Dot(&oc) - radius*radius
	discriminant := b*b - 4*a*c
	return discriminant > 0
}
