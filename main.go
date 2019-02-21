package main

import (
	"bitbucket.org/reidev/restlib/stl"
	"fmt"
	"log"
	"os"
)

//Define the global variables that are setup in the main
//var calcsRepo calcs.Repo

func main() {

	//file, err := os.Open("/Users/mcgurn/Desktop/TorusAscii.stl")
	file, err := os.Open("/Users/mcgurn/Downloads/AirAssis_mod_si.stl")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	stlMesh, err := stl.ReadMesh(file)

	stlMesh.SaveAsPng("/Users/mcgurn/Downloads/img.png")

	fmt.Println(stlMesh)
	fmt.Println(err)

}
