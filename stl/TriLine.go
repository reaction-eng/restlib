// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package stl

type TriLine struct {
	FaceTriEdge [3]Edge
}

func (ele Element) GetEdge() *TriLine {

	return &TriLine{FaceTriEdge: [3]Edge{{starPoint: ele.Nodes[0], endPoint: ele.Nodes[1], count: 0},
		{starPoint: ele.Nodes[1], endPoint: ele.Nodes[2], count: 0},
		{starPoint: ele.Nodes[2], endPoint: ele.Nodes[0], count: 0},
	},
	}
}

// see two edge is differ
func (thisEd *Edge) HasEdge(othTriLine *TriLine) (bool, int8) {

	if thisEd.EdgeDiff(&othTriLine.FaceTriEdge[0]) {
		thisEd.count = 1
		othTriLine.FaceTriEdge[0].count = 1
		return true, 0
	} else if thisEd.EdgeDiff(&othTriLine.FaceTriEdge[1]) {
		othTriLine.FaceTriEdge[1].count = 1
		thisEd.count = 1
		return true, 1
	} else if thisEd.EdgeDiff(&othTriLine.FaceTriEdge[2]) {
		thisEd.count = 1
		othTriLine.FaceTriEdge[2].count = 1
		return true, 2
	} else {
		return false, 3
	} // means that this edge is not included in this othTirLines
}
