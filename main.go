package main

import (
	"bitbucket.org/reidev/restlib/stl"
	"fmt"
	"log"
	"os"
)

//Define the global variables that are setup in the main
//var calcsRepo calcs.Repo

func main2() {

	file, err := os.Open("/Users/mcgurn/Desktop/TorusAscii.stl")
	//file, err := os.Open("/Users/mcgurn/Downloads/AirAssis_mod_si.stl")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	stlMesh, err := stl.ReadMesh(file)

	//Now try writing it
	// Write the body to file
	out, err := os.Create("/Users/mcgurn/Downloads/output.stl")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	//Output the mesh
	stlMesh.WriteMeshBinary(out)

	fmt.Println(stlMesh)
	fmt.Println(err)

}

func main() {

	//Layout some points
	pts := []stl.Vertex{
		{0.0, 0.0, 0.0},
		{.25, .25, 0.0},
		{.3, .5, 0.0},
		{.5, .6, 0.0},
		{.7, .8, 0.0},
		{0.0, 1.0, 0.0},
	}

	mesh, err := stl.RotateAndCreateMesh(pts, 4)

	//Now try writing it
	// Write the body to file
	out, err := os.Create("/Users/mcgurn/Downloads/output.stl")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	//Output the mesh
	mesh.WriteMeshAscii(out)

}
