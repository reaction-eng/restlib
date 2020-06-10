// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/reaction-eng/restlib/notification"
	"github.com/reaction-eng/restlib/stl"
	"github.com/reaction-eng/restlib/users"
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

	greetingNotif := notification.Notification{
		Message:    "Sup my man!",
		Priority:   5,
		Expiration: time.Now(),
		Send:       time.Now(),
		UserID:     1,
	}

	//userGrant := users.BasicUser{
	//	Id_:2,
	//	Email_:"keller@reaction-eng.com",
	//	Token_:"1234567890",
	//}
	//config, _ := configuration.NewConfiguration("config.mysql.json", "config.host.json")
	localSql, err := sql.Open("mysql", "root:P1p3sh0p@tcp(:3306)/localDB?parseTime=true") //"root:P1p3sh0p@tcp(:3306)/localDB?parseTime=true"

	sqlConnectiont, err := users.NewRepoMySql(localSql)

	//include some kind of db call to get email to who'm we send to
	userG, err := sqlConnectiont.GetUser(2)

	if err != nil {
		log.Println("Getting user through error. ", err.Error())
	}

	dumNotifier := notification.NewDummyNotifier()
	webNotif := notification.NewWebPushNotifier("config.deleteME.json", localSql)
	//mailNotifier := Notification.NewEmailNotifier()

	err = dumNotifier.Notify(greetingNotif, userG)
	err = webNotif.Notify(greetingNotif, userG)

}
