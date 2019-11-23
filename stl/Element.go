// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package stl

//Store the element it is a vertex
type Element struct {
	Nodes [3]Vertex
}

//Compute a vertex
func NewElement(v0, v1, v2 Vertex) *Element {

	return &Element{
		Nodes: [3]Vertex{
			v0, v1, v2,
		},
	}
}

//Compute a vertex
func (ele Element) getNormal() Vertex {

	//Get the element sides
	ele1Vec := ele.Nodes[1].minus(&ele.Nodes[0])
	ele2Vec := ele.Nodes[2].minus(&ele.Nodes[1])

	//Cross them
	norm := ele1Vec.cross(&ele2Vec)

	//Take the norm and return
	norm.norm()

	return norm
}

//Compute a vertex
func (ele Element) GetNormalReverse() Vertex {

	//Get the element sides
	ele1Vec := ele.Nodes[0].minus(&ele.Nodes[1])
	ele2Vec := ele.Nodes[2].minus(&ele.Nodes[1])

	//Cross them
	norm := ele1Vec.cross(&ele2Vec)

	//Take the norm and return
	norm.norm()

	return norm
}

//Translation and reverse the normal direction
func (ele Element) Translation(normTrans *Vertex, ExtrudLen float64) *Element {

	return &Element{
		Nodes: [3]Vertex{ele.Nodes[0].trans(normTrans, ExtrudLen),
			ele.Nodes[2].trans(normTrans, ExtrudLen),
			ele.Nodes[1].trans(normTrans, ExtrudLen),
		},
	}
}

//Translation and the same the normal direction
func (ele Element) Move(normTrans *Vertex, ExtrudLen float64) *Element {

	return &Element{
		Nodes: [3]Vertex{ele.Nodes[0].trans(normTrans, ExtrudLen),
			ele.Nodes[1].trans(normTrans, ExtrudLen),
			ele.Nodes[2].trans(normTrans, ExtrudLen),
		},
	}
}
