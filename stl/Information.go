package stl

import "math"

func (mesh *Mesh) GetBounds() (low [3]float64, high [3]float64) {

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
	return

}

func (mesh *Mesh) GetLength() float64 {

	//Get the bounds
	low, high := mesh.GetBounds()

	//Compute the distance
	dist := 0.0
	for dir := 0; dir < len(low); dir++ {
		dist += (high[dir] - low[dir]) * (high[dir] - low[dir])
	}
	dist = math.Sqrt(dist)

	return dist
}
