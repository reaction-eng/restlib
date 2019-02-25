package stl

import (
	"errors"
	"math"
)

func RotateAndCreateMesh(pts []Vertex, slices int) (*Mesh, error) {

	//We need at least three pts
	if len(pts) < 3 {
		return nil, errors.New("insufficient points for generating mesh")
	}

	//Store their index for future use
	strCap := pts[0]
	endCap := pts[len(pts)-1]

	//Now get the body pts
	baseBodyPts := pts[1 : len(pts)-1]

	//Store the bodypts length
	bodyPtLen := len(baseBodyPts)

	//Compute how to rotate it
	deltaOmega := 2.0 * math.Pi / float64(slices)

	//Now we need to build the elements
	elements := make([]Element, 0)

	//Now march over every slice
	for s := 0; s < slices; s++ {
		//compute the start and end omega
		startOmega := float64(s) * deltaOmega
		endOmega := float64(s+1) * deltaOmega

		//Now rotate the pts and add them
		startRot := NewRotationMatrix(strCap, endCap, startOmega)
		endRot := NewRotationMatrix(strCap, endCap, endOmega)

		//Step up each step on the ladder
		for step := 0; step < bodyPtLen-1; step++ {
			//Build the left
			left := NewElement(
				startRot.rotate(&baseBodyPts[step]),
				endRot.rotate(&baseBodyPts[step]),
				startRot.rotate(&baseBodyPts[step+1]),
			)
			right := NewElement(
				endRot.rotate(&baseBodyPts[step]),
				endRot.rotate(&baseBodyPts[step+1]),
				startRot.rotate(&baseBodyPts[step+1]),
			)

			//Add them elements
			elements = append(elements, *left, *right)
		}

		//Add the end pts
		elements = append(elements,
			*NewElement(
				endRot.rotate(&strCap),
				endRot.rotate(&baseBodyPts[0]),
				startRot.rotate(&baseBodyPts[0]),
			),
			*NewElement(
				startRot.rotate(&baseBodyPts[bodyPtLen-1]),
				endRot.rotate(&baseBodyPts[bodyPtLen-1]),
				startRot.rotate(&endCap),
			),
		)

	}

	//Return the stl
	return &Mesh{
		Elements: elements,
	}, nil

}
