// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package google

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

/**
 * An interface type to hold a document or directory
 */
type Item interface {
	GetId() string
	GetName() string
	GetDate() *time.Time
}

//Hold a base type document
type File struct {
	//Keep a boolean if it is a file
	Id string `json:"id"`

	//Hold the item
	Name string `json:"name"`

	//Hold if we should hide the item
	HideListing bool `json:"hideListing"`

	//Keep a date if useful
	Date *time.Time `json:"date"`
}

//Hold a base type document
func (file File) GetId() string {
	return file.Id
}
func (file File) GetName() string {
	return file.Name
}
func (file File) GetDate() *time.Time {
	return file.Date
}

//Hold the structs need to create a tree
type Document struct {
	File

	Type string `json:"type"`

	//Keep the Preview
	Preview string `json:"preview"`

	//Thumbnail Image
	ThumbnailUrl string `json:"thumbnail"`

	//Also Keep the parent id
	ParentId string `json:"parentid"`
}

func (dir Document) MarshalJSON() ([]byte, error) {
	//Store a ddir copy
	type fakeItem Document

	//define the item type
	type ItemWithType struct {
		fakeItem
		InternalItemType string
	}

	return json.Marshal(ItemWithType{
		(fakeItem)(dir),
		"Document",
	})
}

//Hold the struct need to create a tree
type Directory struct {
	File

	//Keep the type for directory anyways
	Type string `json:"type"`

	//Hold a list of items
	Items []Item `json:"items"`

	//Keep the parent id, unless this is root and then it is null
	ParentId string `json:"parentid"`
}

/**
Custom marashaler
*/
func (dir Directory) MarshalJSON() ([]byte, error) {
	//Store a ddir copy
	type fakeItem Directory

	//define the item type
	type ItemWithType struct {
		fakeItem
		InternalItemType string
	}

	return json.Marshal(ItemWithType{
		(fakeItem)(dir),
		"Directory",
	})
}

func (dir *Directory) UnmarshalJSON(b []byte) error {
	// First, deserialize everything into a map of map
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	//Extract a copy of the items
	var rawItemsList []*json.RawMessage
	err = json.Unmarshal(*objMap["items"], &rawItemsList)
	if err != nil {
		return err
	}

	//Now remove the value
	objMap["items"] = nil

	//type def dir
	type dirTmp Directory

	//Now restore back
	bytes, err := json.Marshal(objMap)
	//And back into the object
	err = json.Unmarshal(bytes, (*dirTmp)(dir))

	//Now decode each of the items in the list
	dir.Items = make([]Item, len(rawItemsList))

	var m map[string]interface{}
	for index, rawMessage := range rawItemsList {
		err = json.Unmarshal(*rawMessage, &m)
		if err != nil {
			return err
		}

		//Get the
		InternalItemType := fmt.Sprint(m["InternalItemType"])

		// Depending on the type, we can run json.Unmarshal again on the same byte slice
		// But this time, we'll pass in the appropriate struct instead of a map
		switch InternalItemType {
		case "Document":
			var tmp Document
			err := json.Unmarshal(*rawMessage, &tmp)
			if err != nil {
				return err
			}
			dir.Items[index] = &tmp
			break
		case "Directory":
			var tmp Directory
			err := json.Unmarshal(*rawMessage, &tmp)
			if err != nil {
				return err
			}
			dir.Items[index] = &tmp
			break
		case "Form":
			var tmp Form
			err := json.Unmarshal(*rawMessage, &tmp)
			if err != nil {
				return err
			}
			dir.Items[index] = &tmp
			break
		case "Event":
			var tmp Event
			err := json.Unmarshal(*rawMessage, &tmp)
			if err != nil {
				return err
			}
			dir.Items[index] = &tmp
			break
		default:
			return errors.New("Unsupported type found!")
		}
	}

	return nil
}

//Define a function
type ItemFunc func(item Item)

//Hold the struct need to create a tree
func (dir *Directory) MarchDownTree(doOnFolder bool, do ItemFunc) {
	if dir == nil {
		return
	}

	//Do this to the item
	if doOnFolder {
		do(dir)
	}

	//March over each item in the dir
	for _, item := range dir.Items {
		//If it is a dir, do a recurisive dive
		if asDir, isDir := item.(*Directory); isDir {
			asDir.MarchDownTree(doOnFolder, do)
		} else {
			//Just apply the function
			do(item)
		}

	}

}

//Store the forms metadata
type FormMetaData struct {
	Title string `json:"title"`

	EmailTo []string `json:"emailTo"`

	EmailTemplate string `json:"emailTemplate"`

	DriveInfo []FormDriveInfo `json:"driveInfo"`

	RequiredPerm []string `json:"requiredPerm"`
}

//Store the forms metadata
type FormDriveInfo struct {
	SheetId   string `json:"sheetId"`
	SheetName string `json:"sheetName"`
}

//Hold the struct need to create a tree
type Form struct {
	File
	//Keep the type for directory anyways
	Type string `json:"type"`

	//Keep the parent id, unless this is root and then it is null
	ParentId string `json:"parentid"`

	//Hold the item
	Metadata FormMetaData `json:"metadata"`

	//Hold a list of items
	JSONSchema map[string]interface{} `json:"JSONSchema"`

	//Keep the parent id, unless this is root and then it is null
	UISchema map[string]interface{} `json:"UISchema"`
}

func (dir Form) MarshalJSON() ([]byte, error) {
	//Store a ddir copy
	type fakeItem Form

	//define the item type
	type ItemWithType struct {
		fakeItem
		InternalItemType string
	}

	return json.Marshal(ItemWithType{
		(fakeItem)(dir),
		"Form",
	})
}

//Hold the struct need to create a tree
type Event struct {
	File

	//Keep a boolean if it is a file
	InfoId string `json:"infoId"`

	//Keep a boolean if it is a file
	SignupId string `json:"signupId"`

	//Keep the parent id, unless this is root and then it is null
	ParentId string `json:"parentid"`
}

func (dir Event) MarshalJSON() ([]byte, error) {
	//Store a ddir copy
	type fakeItem Event

	//define the item type
	type ItemWithType struct {
		fakeItem
		InternalItemType string
	}

	return json.Marshal(ItemWithType{
		(fakeItem)(dir),
		"Event",
	})
}

// ByAge implements sort.Interface for []Person based on
// the Age field.
type ByDate []Item

func (a ByDate) Len() int      { return len(a) }
func (a ByDate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool {
	//Get both dates
	iDate := a[i].GetDate()
	jDate := a[j].GetDate()

	//If
	if iDate == nil && jDate == nil {
		return false
	} else if iDate == nil {
		return true
	} else if jDate == nil {
		return false
	} else {
		return (*jDate).Before(*iDate)
	}
}
