// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package users

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/**
Define a struct for Repo for use with users
*/
type RepoMongo struct {
	//Hold on to the sql databased
	db *mongo.Collection

	//Also store the table name
	collectionName string
}

//Provide a metho to make a new Mongo Repo
func NewRepoMongo(locOfDB string, dbName string, collection string) *RepoMongo {
	db := ConnectToDB(locOfDB, dbName)

	coll := db.Collection(collection)

	toRet := RepoMongo{
		db:             coll,
		collectionName: collection,
	}
	return &toRet
}

//Provide a method to make a new UserRepoSql
/**
Get the user with the email.  An error is thrown is not found
*/
func (repo *RepoMongo) GetUserByEmail(email string) (User, error) {
	var user *BasicUser

	//query for email
	cursor := repo.QueryBD(bson.D{})

	//get first result
	user = GetUserFromCursor(*cursor)

	//Store if this is activated
	//user.activated_ = activationDate.Valid
	//user.passwordlogin_ = len(user.password_) > 0

	return user, nil
}

/**
Get the user with the ID.  An error is thrown is not found
*/
//func (repo *RepoMongo)GetUser(id int) (User, error){}
//
///**
//Add User
//*/
//func (repo *RepoMongo)AddUser(user User) (User, error){}
//
///**
//Update User
//*/
//func (repo *RepoMongo)UpdateUser(user User) (User, error){}
//
///**
//Activate User
//*/
//func (repo *RepoMongo)ActivateUser(user User) error{}
//
///**
//Allow databases to be closed
//*/
//func (repo *RepoMongo)CleanUp(){}
//
///**
//Create empty user
//*/
//func (repo *RepoMongo)NewEmptyUser() User{}
//
///**
//List all users
//*/
//func (repo *RepoMongo)ListUsers() ([]int, error){}

//Connect to a db, returns pointer to db
func ConnectToDB(locOfDB string, dbName string) *mongo.Database {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(locOfDB))
	if err != nil {
		log.Println("Couldn't connect to Mongo, ERR:", err)
		return nil
	}

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Println("MongoDB Error, ERR:", err)
		return nil
	}

	log.Println("Connected to Mongo successfully.")
	db := client.Database(dbName)

	return db
}

//query DB, returns cursor
func (repo *RepoMongo) QueryBD(queryString bson.D) *mongo.Cursor {
	ctx := context.Background()

	cursor, err := repo.db.Find(ctx, queryString, options.Find().SetProjection(bson.D{{"_id", 0}}))
	defer cursor.Close(ctx) //tell it to close when done.

	if err != nil {
		log.Println("MongoProcessor GetAll error ->", err)
		return nil
	}

	return cursor
}

//
func GetUserFromCursor(cursor mongo.Cursor) *BasicUser {
	cursor.Next(context.Background())
	jsonElem := bson.M{}

	if err := cursor.Decode(&jsonElem); err != nil {
		log.Println("ERROR from mongo cursor in GetUserFromCursor ->", err)
		return nil
	}

	basicUser := BasicUser{
		Id_:            int(jsonElem["ID"].(float64)),
		Email_:         jsonElem["Email"].(string),
		password_:      jsonElem["Password"].(string),
		Token_:         jsonElem["Token"].(string),
		activated_:     jsonElem["ID"].(bool),
		passwordlogin_: jsonElem["ID"].(bool),
	}

	return &basicUser
}
