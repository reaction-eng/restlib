// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package stl

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/reaction-eng/restlib/utils"
	"io"
	"regexp"
	"strconv"
	"strings"
)

//Store the max element size
const maxNumEle = 1000000

//Hold everything in the mesh
type Mesh struct {
	//Store a list of Elements
	Elements []Element

	//Store the name of the info
	Name string
}

//Read the mesh from a file
func ReadMesh(in io.Reader) (*Mesh, error) {

	//Check for nil in
	if in == nil {
		return nil, errors.New("invalid reader")
	}

	//Wrap the reader in a buf io
	bufIn := bufio.NewReader(in)

	//Assume we are bool
	binary := true

	//Read in 80 bytes
	stringData, err := bufIn.Peek(80)

	//If there is an error
	if err != nil {
		return nil, err
	}

	//See if it is a string
	string := string(stringData)

	//See if it binary
	if strings.HasPrefix(strings.ToLower(string), "solid") {
		binary = false
	}

	if binary {
		return readMeshFromBinary(bufIn)

	} else {
		return readMeshFromAscii(bufIn)
	}

}

//Hold everything in the mesh
func readMeshFromBinary(in *bufio.Reader) (*Mesh, error) {
	//Create a new mesh with elements
	mesh := &Mesh{}

	//REad the first 80 bytes
	titleBytes := make([]byte, 80)

	//Read it
	_, err := in.Read(titleBytes)

	//If there was a problem return
	if err != nil {
		return nil, err
	}

	//Remove any invalid info from the string
	titleBytes = bytes.Trim(titleBytes, "\x00")

	//Store the name
	mesh.Name = string(titleBytes)

	//Now determine the number of elements
	var numEle int32
	err = binary.Read(in, binary.LittleEndian, &numEle)
	if err != nil {
		return nil, err
	}
	if numEle > maxNumEle {
		return nil, errors.New("invalid stl or more than allowed elements: " + strconv.Itoa(int(numEle)))
	}

	//Size the array
	mesh.Elements = make([]Element, numEle)

	//For each ele
	for ele := int32(0); ele < numEle; ele++ {
		//Read the normal vec
		var normVec Vertex

		//Load it
		err = binary.Read(in, binary.LittleEndian, &normVec)
		if err != nil {
			return nil, err
		}

		//Build an element
		newEle := Element{}

		//Read each vertexInEle
		for v := 0; v < len(newEle.Nodes); v++ {
			//Read it
			err = binary.Read(in, binary.LittleEndian, &newEle.Nodes[v])
			if err != nil {
				return nil, err
			}

		}

		//Add in the element
		mesh.Elements[ele] = newEle

		//stls allows you specify other info in attribute byte count
		var attributeByteCount uint16
		err = binary.Read(in, binary.LittleEndian, &attributeByteCount)
		if err != nil {
			return nil, err
		}

		//We don't need the attribute right now so just read advance by that many bytes
		_, err = in.Discard(int(attributeByteCount))
		if err != nil {
			return nil, err
		}
	}

	return mesh, nil

}

//Read in a binary vect

//Hold everything in the mesh
func readMeshFromAscii(in *bufio.Reader) (*Mesh, error) {
	//Create a new mesh with elements
	mesh := &Mesh{
		Elements: make([]Element, 0),
	}

	//Convert to a scanner
	scanner := bufio.NewScanner(in)

	//Scan the first line
	scanner.Scan()
	mesh.Name = scanner.Text()

	//march over each line
	for scanner.Scan() {
		//Get the line
		line := scanner.Text()

		facetLine := splitLine(line)
		firstArg := strings.TrimSpace(facetLine[0])

		//If this is a facet
		if firstArg == ("facet") {
			//Build the norm
			//norm, err := buildVertex(facetLine[0 + 2], facetLine[1 + 2], facetLine[2 + 2])

			//Skip the line
			scanner.Scan()

			//Build an element
			ele := Element{}

			//Now over each vertex
			for i := 0; i < len(ele.Nodes); i++ {
				//Get the node line
				scanner.Scan()
				nodeLine := splitLine(scanner.Text())

				//Store the vertex
				node, err := buildVertex(nodeLine[0+1], nodeLine[1+1], nodeLine[2+1])

				if err != nil {
					return nil, err
				}

				//Store the node
				ele.Nodes[i] = node

			}

			//Add in the element
			mesh.Elements = append(mesh.Elements, ele)

		}

	}

	//Check for error
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return mesh, nil

}

//Build the vertex from a line
func buildVertex(v0 string, v1 string, v2 string) (Vertex, error) {
	//Build a new vertex

	//Convert each one to float64
	v0Float, err := strconv.ParseFloat(v0, 32)
	if err != nil {
		return Vertex{}, err
	}
	v1Float, err := strconv.ParseFloat(v1, 32)
	if err != nil {
		return Vertex{}, err
	}
	v2Float, err := strconv.ParseFloat(v2, 32)
	if err != nil {
		return Vertex{}, err
	}
	return Vertex{float32(v0Float), float32(v1Float), float32(v2Float)}, nil
}

//Build the vertex from a line
func splitLine(line string) []string {
	//Build a new vertex
	//Clean up the line
	line = strings.TrimSpace(line)
	whtSpace := regexp.MustCompile(`\s+`)
	line = whtSpace.ReplaceAllString(line, ",")

	//Now split and get the first argument
	return strings.Split(line, ",")
}

//Returns a subset of elements
func (mesh *Mesh) GetSubsetMesh(ints []int) (*Mesh, error) {
	//Create a new mesh with elements
	newMesh := &Mesh{
		Elements: make([]Element, 0),
	}

	//Now copy over the ints
	for _, index := range ints {
		//Make sure it is valid
		if index < 0 || index > len(mesh.Elements) {
			return nil, errors.New("invalid index for element")
		}

		//Ok add it
		newMesh.Elements = append(newMesh.Elements, mesh.Elements[index])

	}
	return newMesh, nil

}

func (mesh *Mesh) ConvertToSI(unit utils.Unit) {
	//Get the si factor as a float32
	factor := float32(unit.GetFactorSI())

	//Update all of the positions to convert to si
	//March over each node
	for e := range mesh.Elements {
		//For each vertex
		for n := range mesh.Elements[e].Nodes {
			//Add to the value
			for dir := 0; dir < len(mesh.Elements[e].Nodes); dir++ {
				mesh.Elements[e].Nodes[n][dir] *= factor
			}

		}
	}

}
