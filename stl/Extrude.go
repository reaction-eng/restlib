// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package stl

type ExtrudeFace struct {
	//Store the matrix
	matrix [3][3]float64 //[row][col]

	//And the start and end pts
	start       Vertex
	end         Vertex
	startStrech Vertex
	endStrecth  Vertex
}

func Extrude(start Vertex, end Vertex, normExtrude *Vertex, ExtrudLen float64) *ExtrudeFace {
	//Create the struct
	ExFaceMatrix := &ExtrudeFace{
		start: start,
		end:   end,
	}
	// check if normExtrude is unit one
	ExFaceMatrix.startStrech = start.trans(normExtrude, ExtrudLen)
	ExFaceMatrix.endStrecth = end.trans(normExtrude, ExtrudLen)
	//ExFaceMatrix.startStrech[0] = start[0]+normExtrude[0]*ExtrudLen
	//ExFaceMatrix.startStrech[1] = start[1]+normExtrude[1]*ExtrudLen
	//ExFaceMatrix.startStrech[2] = start[2]+normExtrude[2]*ExtrudLen
	//
	//
	//ExFaceMatrix.endStrecth[0] = end[0]+normExtrude[0]*ExtrudLen
	//ExFaceMatrix.endStrecth[1] = end[1]+normExtrude[1]*ExtrudLen
	//ExFaceMatrix.endStrecth[2] = end[2]+normExtrude[2]*ExtrudLen

	return ExFaceMatrix

}
