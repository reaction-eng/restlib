// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package stl

import "math"

type RotationMatrix struct {
	//Store the matrix
	matrix [3][3]float32 //[row][col]

	//And the start and end pts
	start   Vertex
	end     Vertex
	lineVec Vertex
	lineMag float32
}

//Create the rotation matrix
func NewRotationMatrix(start Vertex, end Vertex, theta float64) *RotationMatrix {
	//Create the struct
	rotMat := &RotationMatrix{
		start:   start,
		end:     end,
		lineVec: end.minus(&start),
	}

	//Store the line mag
	rotMat.lineMag = rotMat.lineVec.mag()

	//Compute a unit dir vector
	uv := end.minus(&start)
	uv.norm()

	//Precompute some values
	cos := float32(math.Cos(theta))
	sin := float32(math.Sin(theta))
	oneMinusCos := float32(1.0 - cos)

	//Now build the rotation vector
	rotMat.matrix[0][0] = cos + uv[0]*uv[0]*oneMinusCos
	rotMat.matrix[0][1] = uv[0]*uv[1]*oneMinusCos - uv[2]*sin
	rotMat.matrix[0][2] = uv[0]*uv[2]*oneMinusCos + uv[1]*sin

	rotMat.matrix[1][0] = uv[1]*uv[0]*oneMinusCos + uv[2]*sin
	rotMat.matrix[1][1] = cos + uv[1]*uv[1]*oneMinusCos
	rotMat.matrix[1][2] = uv[1]*uv[2]*oneMinusCos - uv[0]*sin

	rotMat.matrix[2][0] = uv[2]*uv[0]*oneMinusCos - uv[1]*sin
	rotMat.matrix[2][1] = uv[2]*uv[1]*oneMinusCos + uv[0]*sin
	rotMat.matrix[2][2] = cos + uv[2]*uv[2]*oneMinusCos

	return rotMat
}

//Gets the point origin on the line
func (rot *RotationMatrix) getPtOrg(pt *Vertex) Vertex {

	// compute the distance along the line
	stpt := rot.start.minus(pt)
	fracTop := rot.lineVec.dot(&stpt)
	fracBot := rot.lineMag * rot.lineMag

	//Compute the fraction t
	t := -fracTop / (fracBot + 1e-30)

	//Now compute the location
	if t < 0 {
		return rot.start
	} else if t < 1.0 {
		return Vertex{
			t*rot.lineVec[0] + rot.start[0],
			t*rot.lineVec[1] + rot.start[1],
			t*rot.lineVec[2] + rot.start[2],
		}

	} else {
		return rot.end
	}

}

//Create the rotation matrix
func (rot *RotationMatrix) rotate(pt *Vertex) Vertex {
	//Get the org for this pt
	ptOrg := rot.getPtOrg(pt)

	//Now get this pt relative to the pt org
	relPt := pt.minus(&ptOrg)

	//Now rotate it
	var rotatPt Vertex

	for i := 0; i < len(rot.matrix); i++ {
		for j := 0; j < len(rot.matrix[i]); j++ {
			rotatPt[i] += rot.matrix[i][j] * relPt[j]
		}
	}

	//Now add the ptOrg back to the point
	rotatPt.addTo(&ptOrg)

	return rotatPt

}
