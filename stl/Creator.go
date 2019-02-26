package stl

import (
	"errors"
	"math"
)

//Gets the point origin on the line
func getPtAlongLine(str, end *Vertex, frac float32) Vertex {

	// compute the distance along the line
	return Vertex{
		str[0] + (end[0]-str[0])*frac,
		str[1] + (end[1]-str[1])*frac,
		str[2] + (end[2]-str[2])*frac,
	}

}

func IncreaseLineResolution(pts []Vertex, radialRes, longitudinalRes float32) ([]Vertex, error) {

	//We need at least three pts
	if len(pts) < 2 {
		return nil, errors.New("insufficient points for generating mesh")
	}

	//Make a new list of pts
	newList := make([]Vertex, 0)

	//compute the unit vector for the vert direction
	vertUnit := pts[len(pts)-1].minus(&pts[0])
	vertUnit.norm()

	//Now step over each line
	for l := 0; l < len(pts)-1; l++ {
		//Get the start and end pts
		strPt := pts[l]
		endPt := pts[l+1]

		//Get the line vector
		linVec := endPt.minus(&strPt)

		//Determine the number of sub divisions
		vertComp := linVec.dot(&vertUnit)
		//Compute the size of the vector in the vert dir
		vertCompDir := vertUnit.scaleCopy(vertComp)

		//The remainder is the rad comp
		radCompVect := vertCompDir.minus(&linVec)
		radComp := radCompVect.mag()

		//Setup the number of sublines
		numSubLines := 1

		//Compute the req in each dir
		vertSubLines := int(vertComp / longitudinalRes)
		radSubLines := int(radComp / radialRes)
		if vertSubLines > numSubLines {
			numSubLines = vertSubLines
		}
		if radSubLines > numSubLines {
			numSubLines = radSubLines
		}

		//Add the init line
		newList = append(newList, strPt)

		//Now sub divide by the number of sublines
		for subLine := 1; subLine < numSubLines; subLine++ {
			//Compute the fraction
			frac := float32(subLine) / float32(numSubLines)

			//Compute and add the pt
			newList = append(newList, getPtAlongLine(&strPt, &endPt, frac))

		} //Note the last pt is covered by the first point in the line

	}
	//The last point is not covered with the lines, so add pts
	newList = append(newList, pts[len(pts)-1])

	return newList, nil

}

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
