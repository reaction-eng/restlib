package stl

import (
	"github.com/fogleman/fauxgl"
	"github.com/nfnt/resize"
	"image"
	"image/png"
	"math"
	"os"
)

const (
	scale   = 1    // optional supersampling
	width   = 1920 // output width in pixels
	height  = 1080 // output height in pixels
	fovy    = 30   // vertical field of view in degrees
	near    = 1    // near clipping plane
	eyeFact = .5
	farFac  = 4 //Distance factor for clipping
)

var (
	up    = fauxgl.V(0, 0, 1)                    // up vector
	light = fauxgl.V(-0.75, 1, 0.25).Normalize() // light direction
	color = fauxgl.HexColor("#468966")           // object color
)

func convertToImgVec(vertex *Vertex) fauxgl.Vector {
	return fauxgl.V(float64(vertex[0]), float64(vertex[1]), float64(vertex[2]))
}

func (mesh *Mesh) SaveAsPng(file string) error {

	//Convert to an img
	img, err := mesh.GetImage()

	// Write the body to file
	out, err := os.Create(file)
	if err != nil {
		return err
	}
	defer out.Close()
	return png.Encode(out, img)

}

func (mesh *Mesh) GetImage() (image.Image, error) {

	//Img elements
	imgTriangles := make([]*fauxgl.Triangle, len(mesh.Elements))

	//Copy over each one
	for e := 0; e < len(imgTriangles); e++ {
		//Get the triangle
		ele := mesh.Elements[e]

		//Convert to a tri angle
		imgTriangles[e] = fauxgl.NewTriangleForPoints(
			convertToImgVec(&ele.Nodes[0]),
			convertToImgVec(&ele.Nodes[1]),
			convertToImgVec(&ele.Nodes[2]),
		)

	}

	//Convert the mesh elements to the image elements
	imageMesh := fauxgl.NewTriangleMesh(imgTriangles)

	// fit mesh in a bi-unit cube centered at the origin
	imageMesh.BiUnitCube()

	// smooth the normals
	imageMesh.SmoothNormalsThreshold(fauxgl.Radians(30))

	// create a rendering context
	context := fauxgl.NewContext(width*scale, height*scale)
	context.ClearColorBufferWith(fauxgl.HexColor("#FFF8E3"))

	// create transformation matrix and light direction
	aspect := float64(width) / float64(height)

	//Compute the locations
	centerVertex := mesh.GetCenter()
	center := convertToImgVec(&centerVertex)

	//Make a copy
	eye := center
	objLength := mesh.GetLength()
	eye.X += objLength * eyeFact
	eye.Y += objLength * eyeFact
	eye.Z += objLength * eyeFact

	//Assume the far is
	far := objLength * eyeFact * farFac

	//Get the directon vector
	matrix := fauxgl.LookAt(eye, center, up).Perspective(fovy, aspect, near, far)

	// use bu iltin phong shader
	shader := fauxgl.NewPhongShader(matrix, light, eye)
	shader.ObjectColor = color
	context.Shader = shader

	// render
	context.DrawMesh(imageMesh)

	// downsample image for antialiasing
	image := context.Image()
	image = resize.Resize(width, height, image, resize.Bilinear)

	return image, nil

}

func (mesh *Mesh) GetCenter() Vertex {

	////Get the center
	//var center Vertex
	//
	////March over each node
	//for _, ele := range mesh.Elements {
	//	//For each vertex
	//	for _, node := range ele.Nodes {
	//		//Add to the value
	//		for dir := 0; dir < len(node); dir++ {
	//			center[dir] += node[dir]
	//		}
	//
	//	}
	//}
	//
	////Now take the average
	//for dir := 0; dir < len(center); dir++ {
	//	center[dir] /= float32(len(mesh.Elements) * len(center))
	//}
	//
	//return center

	//Get the center
	var low [3]float64
	var high [3]float64

	//Set default values
	for dir := 0; dir < len(low); dir++ {
		low[dir] = math.MaxFloat64
		high[dir] = -math.MaxFloat64

	}

	//March over each node
	for _, ele := range mesh.Elements {
		//For each vertex
		for _, node := range ele.Nodes {
			//Add to the value
			for dir := 0; dir < len(node); dir++ {
				low[dir] = math.Min(low[dir], float64(node[dir]))
				high[dir] = math.Max(high[dir], float64(node[dir]))

			}

		}
	}

	//Compute the distance
	var center Vertex
	for dir := 0; dir < len(low); dir++ {
		center[dir] = float32(0.5 * (low[dir] + high[dir]))
	}

	return center
}
func (mesh *Mesh) GetLength() float64 {

	//Get the center
	var low [3]float64
	var high [3]float64

	//Set default values
	for dir := 0; dir < len(low); dir++ {
		low[dir] = math.MaxFloat64
		high[dir] = -math.MaxFloat64

	}

	//March over each node
	for _, ele := range mesh.Elements {
		//For each vertex
		for _, node := range ele.Nodes {
			//Add to the value
			for dir := 0; dir < len(node); dir++ {
				low[dir] = math.Min(low[dir], float64(node[dir]))
				high[dir] = math.Max(high[dir], float64(node[dir]))

			}

		}
	}

	//Compute the distance
	dist := 0.0
	for dir := 0; dir < len(low); dir++ {
		dist += (high[dir] - low[dir]) * (high[dir] - low[dir])
	}
	dist = math.Sqrt(dist)

	return dist
}
