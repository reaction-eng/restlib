// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package stl

// Get the corresponding edge for the Nodes
type Edge struct {
	starPoint Vertex
	endPoint  Vertex
	count     int16
}

// see two edge is differ
func (thisEd *Edge) EdgeDiff(othEd *Edge) bool {
	sameDir := (thisEd.starPoint.diff(&othEd.starPoint) && thisEd.endPoint.diff(&othEd.endPoint))
	reverDir := (thisEd.starPoint.diff(&othEd.endPoint) && thisEd.endPoint.diff(&othEd.starPoint))
	return (sameDir || reverDir)
}
