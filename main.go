package main

import (
	"fmt"
	"os"
	"time"

	. "github.com/fogleman/fauxgl"
	"github.com/nfnt/resize"
)

const (
	scale       = 1    // optional supersampling
	width       = 1920 // output width in pixels
	height      = 1080 // output height in pixels
	fovy        = 30   // vertical field of view in degrees
	near        = 1    // near clipping plane
	far         = 10   // far clipping plane
	COLOR_BLACK = "#000000"
	COLOR_WHITE = "#FFFFFF"
	COLOR_GRAY  = "#C9C9C9"
)

var (
	eye    = V(2.5, -1.5, 2.0)             // camera position
	center = V(0.25, 0, 0)                 // view center position
	up     = V(0, 0, 1)                    // up vector
	light  = V(2.0, -2.0, 1.5).Normalize() // light direction
	color  = HexColor(COLOR_GRAY)          // object color
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide the path to the STL file as a command line argument.")
		os.Exit(1)
	}

	startTime := time.Now()

	stlPath := os.Args[1]

	// load a mesh
	mesh, err := LoadSTL(stlPath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Loaded mesh with %d triangles\n", len(mesh.Triangles))

	// fit mesh in a bi-unit cube centered at the origin
	mesh.BiUnitCube()
	// mesh.Transform(Scale(V(0.5, 0.5, 0.5)))

	// smooth the normals
	mesh.SmoothNormalsThreshold(Radians(30))

	// create a rendering context
	context := NewContext(width*scale, height*scale)
	context.ClearColorBufferWith(HexColor(COLOR_WHITE))
	// context.Wireframe = true

	// create transformation matrix and light direction
	aspect := float64(width) / float64(height)
	matrix := LookAt(eye, center, up).Perspective(fovy, aspect, near, far)

	// use builtin phong shader
	shader := NewPhongShader(matrix, light, eye)
	shader.ObjectColor = color
	context.Shader = shader

	// render
	context.DrawMesh(mesh)

	// downsample image for antialiasing
	image := context.Image()
	image = resize.Resize(width, height, image, resize.Bilinear)

	// save image
	SavePNG("out.png", image)

	elapsed := time.Since(startTime)
	fmt.Printf("Image generated in %s\n", elapsed)
}
