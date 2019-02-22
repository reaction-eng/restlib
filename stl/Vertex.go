package stl

import "math"

//Setup the vertex type
type Vertex [3]float32

//Function to subtract
func (v *Vertex) subtract(other *Vertex, result *Vertex) {
	//Simple sub
	result[0] = v[0] - other[0]
	result[1] = v[1] - other[1]
	result[2] = v[2] - other[2]

}

//Function to subtract
func (v *Vertex) minus(other *Vertex) Vertex {
	//New vertex
	var result Vertex

	//Copy the data
	v.subtract(other, &result)

	//Return the result
	return result

}

//Function to cross
func (v *Vertex) cross(oth *Vertex) Vertex {
	//New vertex
	return Vertex{
		v[1]*oth[2] - oth[1]*v[2],
		oth[0]*v[2] - v[0]*oth[2],
		v[0]*oth[1] - oth[0]*v[1],
	}

}

//Function to cross
func (v *Vertex) dot(oth *Vertex) float32 {
	//New vertex
	return v[0]*oth[0] + v[1]*oth[1] + v[2]*oth[2]

}

//Function to cross
func (v *Vertex) norm() {
	//Get the mag
	mag := float32(math.Sqrt(float64(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])))

	//New vertex
	v[0] /= mag + 1E-30
	v[1] /= mag + 1E-30
	v[2] /= mag + 1E-30

}

func (v *Vertex) mag() float32 {
	return float32(math.Sqrt(float64(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])))

}

func (v *Vertex) addTo(other *Vertex) {
	v[0] += other[0]
	v[1] += other[1]
	v[2] += other[2]

}
