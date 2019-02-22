package stl

//Store the element it is a vertex
type Element struct {
	Nodes [3]Vertex
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
