// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

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

func ExtrudeTriAndCreateMesh(pts []Vertex, normExtrude *Vertex, ExtrudLen float64) (*Mesh, error) {

	//We need at least three pts
	if len(pts) < 3 {
		return nil, errors.New("insufficient points for generating mesh")
	}

	//Store their index for future use
	OneCap := pts[0]
	TwoCap := pts[1]
	ThirCap := pts[2]

	var OneStretch Vertex
	var TwoStretch Vertex
	var ThirStretch Vertex

	//
	OneStretch = OneCap.trans(normExtrude, ExtrudLen)
	TwoStretch = TwoCap.trans(normExtrude, ExtrudLen)
	ThirStretch = ThirCap.trans(normExtrude, ExtrudLen)
	//
	//OneStretch[0] = OneCap[0]+ normExtrude[0]*ExtrudLen
	//OneStretch[1] = OneCap[1]+ normExtrude[1]*ExtrudLen
	//OneStretch[2] = OneCap[2]+ normExtrude[2]*ExtrudLen
	//
	//TwoStretch[0] = TwoCap[0]+ normExtrude[0]*ExtrudLen
	//TwoStretch[1] = TwoCap[1]+ normExtrude[1]*ExtrudLen
	//TwoStretch[2] = TwoCap[2]+ normExtrude[2]*ExtrudLen
	//
	//ThirStretch[0] = ThirCap[0]+ normExtrude[0]*ExtrudLen
	//ThirStretch[1] = ThirCap[1]+ normExtrude[1]*ExtrudLen
	//ThirStretch[2] = ThirCap[2]+ normExtrude[2]*ExtrudLen

	////Now get the body pts
	//baseBodyPts := pts[1 : len(pts)-1]
	//
	////Store the bodypts length
	//bodyPtLen := len(baseBodyPts)

	//Now we need to build the elements
	elements := make([]Element, 0)

	//Add the top face
	elements = append(elements,
		*NewElement(
			OneCap,
			TwoCap,
			ThirCap,
		),
	)

	for lineIndex := 0; lineIndex < len(pts)-1; lineIndex++ {
		OneFace := Extrude(pts[lineIndex], pts[lineIndex+1], normExtrude, ExtrudLen)
		left := NewElement(
			OneFace.end,
			OneFace.start,
			OneFace.startStrech,
		)
		right := NewElement(
			OneFace.end,
			OneFace.startStrech,
			OneFace.endStrecth,
		)

		elements = append(elements, *left, *right)

	}
	//Add the face based from endpoint to staring point
	OneFace := Extrude(pts[len(pts)-1], pts[0], normExtrude, ExtrudLen)
	leftEndLoop := NewElement(
		OneFace.end,
		OneFace.start,
		OneFace.startStrech,
	)
	rightEndLoop := NewElement(
		OneFace.end,
		OneFace.startStrech,
		OneFace.endStrecth,
	)
	elements = append(elements, *leftEndLoop, *rightEndLoop)

	// Add the parallel top hat
	elements = append(elements,
		*NewElement(
			OneStretch,
			ThirStretch,
			TwoStretch,
		),
	)

	//Return the stl
	return &Mesh{
		Elements: elements,
	}, nil

}

func (facemesh *Mesh) ExtrudeFaceAndCreateMesh(normExtrude *Vertex, ExtrudLen float64) (*Mesh, error) {

	//We need at least three pts
	if len(facemesh.Elements) < 1 {
		return nil, errors.New(" NOT Empty face for generating mesh")
	}
	bottomFace := make([]Element, 0)

	faceEdge := make([]TriLine, 0)

	for _, ele := range facemesh.Elements {
		// creat the corresponding triLine object to store the corresponding edge of triangular
		faceEdge = append(faceEdge, *ele.GetEdge())

		bottomFace = append(bottomFace, *ele.Translation(normExtrude, ExtrudLen))
		// translation the top face to bottom face, also reverse direciton
	} // the index of faceEdge is same as the element of face mesh

	// storing the edge of face
	boundEdge := make([]Edge, 0)
	for ind, TriangEle := range faceEdge {

		for indSec := ind + 1; indSec < len(faceEdge); indSec++ {

			if TriangEle.FaceTriEdge[0].count == 0 {
				TriangEle.FaceTriEdge[0].HasEdge(&faceEdge[indSec])
			}
			if TriangEle.FaceTriEdge[1].count == 0 {
				TriangEle.FaceTriEdge[1].HasEdge(&faceEdge[indSec])
			}
			if TriangEle.FaceTriEdge[2].count == 0 {
				TriangEle.FaceTriEdge[2].HasEdge(&faceEdge[indSec])
			}
		} // for each face edge search

		// after loop, add the edge into boundEdge that isnot included in other faceEdge
		if TriangEle.FaceTriEdge[0].count == 0 {
			boundEdge = append(boundEdge, TriangEle.FaceTriEdge[0])
		}
		if TriangEle.FaceTriEdge[1].count == 0 {
			boundEdge = append(boundEdge, TriangEle.FaceTriEdge[1])
		}
		if TriangEle.FaceTriEdge[2].count == 0 {
			boundEdge = append(boundEdge, TriangEle.FaceTriEdge[2])
		}

	} // end of total search

	//based on the bound edge and create the corrponding face/elemnt
	for _, line := range boundEdge {
		OneFace := Extrude(line.starPoint, line.endPoint, normExtrude, ExtrudLen)
		left := NewElement(
			OneFace.end,
			OneFace.start,
			OneFace.startStrech,
		)
		right := NewElement(
			OneFace.end,
			OneFace.startStrech,
			OneFace.endStrecth,
		)

		facemesh.Elements = append(facemesh.Elements, *left, *right)
	}
	// Add the parallel bottom surface
	facemesh.Elements = append(facemesh.Elements, bottomFace...)

	//Return the stl
	return facemesh, nil

}
