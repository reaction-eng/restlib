package stl

import (
	"fmt"
	"math"
)

func RotateAndCreateMesh(pts []Vertex) {

	//Output each pt
	fmt.Println("x y z data")
	for _, pt := range pts {
		fmt.Println(fmt.Sprint(pt[0]) + " " + fmt.Sprint(pt[1]) + " " + fmt.Sprint(pt[2]) + " 0.0")
	}

	fmt.Println("x y z data")

	//Rotate around
	const slices = 10
	const deltaTheta = 2.0 * math.Pi / slices

	for s := 0; s < slices; s++ {
		theta := deltaTheta * float64(s)

		//Build a rotation matrix
		rotMat := NewRotationMatrix(pts[0], pts[len(pts)-1], theta)

		//Rotate this pt
		rotatedPts := make([]Vertex, 0)

		//Now get each pt
		for _, pt := range pts[1 : len(pts)-1] {
			rotatedPts = append(rotatedPts, rotMat.rotate(&pt))
		}

		for _, pt := range rotatedPts {
			fmt.Println(fmt.Sprint(pt[0]) + " " + fmt.Sprint(pt[1]) + " " + fmt.Sprint(pt[2]) + " " + fmt.Sprint(theta))
		}

	}

}
