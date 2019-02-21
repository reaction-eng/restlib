package main

import (
	"bitbucket.org/reidev/restlib/stl"
	"fmt"
	"image/png"
	"log"
	"os"
)

//Define the global variables that are setup in the main
//var calcsRepo calcs.Repo

func main() {

	//file, err := os.Open("/Users/mcgurn/Desktop/TorusAscii.stl")
	file, err := os.Open("/Users/mcgurn/Desktop/Torus.stl")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	stlMesh, err := stl.ReadMesh(file)

	fmt.Println(stlMesh)
	fmt.Println(err)

	//Convert to an img
	img, err := stlMesh.GetImage()

	// Write the body to file
	out, err := os.Create("/Users/mcgurn/Downloads/img.png")
	if err != nil {
		log.Fatal(err)
	}
	err = png.Encode(out, img)
	if err != nil {
		log.Fatal(err)
	}
}
