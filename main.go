package main

import (
	"bitbucket.org/reidev/restlib/Notification"
	"bitbucket.org/reidev/restlib/stl"
	"fmt"
	"log"
	"os"
	"time"
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

func main3() {

	//Layout some points
	pts := []stl.Vertex{
		{0.0, 0.0, 0.0},
		{0.5, 0.5, 0.0},
		{.75, 0.75, 0.0},
		{.3, 0.9, 0.0},
		{.6, 1.3, 0.0},
		{.9, 1.3, 0.0},
		{.9, 1.6, 0.0},
		{0.0, 1.6, 0.0},
	}

	incPts, _ := stl.IncreaseLineResolution(pts, 0.1, 1.0)
	//incPts := pts
	mesh, _ := stl.RotateAndCreateMesh(incPts, 10)

	//Now try writing it
	// Write the body to file
	out, err := os.Create("/Users/mcgurn/Downloads/output.stl")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	//Now try writing it
	// Write the body to file
	ptsout, err := os.Create("/Users/mcgurn/Downloads/output.pts")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	//Now try writing it
	// Write the body to file
	triout, err := os.Create("/Users/mcgurn/Downloads/output.tri")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	//Output the mesh
	mesh.WriteMeshAscii(out)
	mesh.WriteUintahPts(ptsout)
	mesh.WriteUintahTri(triout)

}

func main() {

	greetingNotif := Notification.Notification{
		Message:    "Sup my man!",
		Priority:   5,
		Expiration: time.Now(),
		Send:       time.Now(),
		User:       "It's G",
	}

	dumNotifier := Notification.NewDummyNotifier()
	mailNotifier := Notification.NewSlackNotifier()

	err := dumNotifier.Notify(greetingNotif)
	if err != nil {
		log.Println("Somehow got an error ->", err.Error())
	}

	err = mailNotifier.Notify(greetingNotif)

}
