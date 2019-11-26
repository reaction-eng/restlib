// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package google

import (
	"encoding/json"
	"fmt"

	"github.com/reaction-eng/restlib/file"
)

type gDirectory struct {
	gFile

	//Keep the type for directory anyways
	Type string `json:"type"`

	//Hold a list of Items
	Items []file.Item `json:"Items"`

	//Keep the parent Id, unless this is root and then it is null
	ParentId string `json:"parentid"`
}

func (dir *gDirectory) GetType() string {
	return dir.Type
}

func (dir *gDirectory) GetItems() []file.Item {
	return dir.Items
}

func (dir *gDirectory) GetParentId() string {
	return dir.ParentId
}

func (dir *gDirectory) ForEach(doOnFolder bool, do file.ItemFunc) {
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
		if asDir, isDir := item.(*gDirectory); isDir {
			asDir.ForEach(doOnFolder, do)
		} else {
			//Just apply the function
			do(item)
		}
	}
}

/**
Custom marashaler
*/
func (dir gDirectory) MarshalJSON() ([]byte, error) {
	//Store a ddir copy
	type fakeItem gDirectory

	//define the item type
	type ItemWithType struct {
		fakeItem
		InternalItemType string
	}

	return json.Marshal(ItemWithType{
		(fakeItem)(dir),
		"gDirectory",
	})
}

func (dir *gDirectory) UnmarshalJSON(b []byte) error {
	// First, deserialize everything into a map of map
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	//Extract a copy of the Items
	var rawItemsList []*json.RawMessage
	err = json.Unmarshal(*objMap["Items"], &rawItemsList)
	if err != nil {
		return err
	}

	//Now remove the value
	objMap["Items"] = nil

	//type def dir
	type dirTmp gDirectory

	//Now restore back
	bytes, err := json.Marshal(objMap)
	//And back into the object
	err = json.Unmarshal(bytes, (*dirTmp)(dir))

	//Now decode each of the Items in the list
	dir.Items = make([]file.Item, len(rawItemsList))

	for index, rawMessage := range rawItemsList {
		var m map[string]interface{}

		err = json.Unmarshal(*rawMessage, &m)
		if err != nil {
			return err
		}

		//Get the
		InternalItemType := fmt.Sprint(m["InternalItemType"])

		// Depending on the type, we can run json.Unmarshal again on the same byte slice
		// But this time, we'll pass in the appropriate struct instead of a map
		switch InternalItemType {
		case "gDocument":
			var tmp gDocument
			err := json.Unmarshal(*rawMessage, &tmp)
			if err != nil {
				return err
			}
			dir.Items[index] = &tmp
			break
		case "gDirectory":
			var tmp gDirectory
			err := json.Unmarshal(*rawMessage, &tmp)
			if err != nil {
				return err
			}
			dir.Items[index] = &tmp
			break
		case "gForm":
			var tmp gForm
			err := json.Unmarshal(*rawMessage, &tmp)
			if err != nil {
				return err
			}
			dir.Items[index] = &tmp
			break
		case "gEvent":
			var tmp gEvent
			err := json.Unmarshal(*rawMessage, &tmp)
			if err != nil {
				return err
			}
			dir.Items[index] = &tmp
			break
		default:
			var tmp gFile
			err := json.Unmarshal(*rawMessage, &tmp)
			if err != nil {
				return err
			}
			dir.Items[index] = &tmp
			break
		}
	}

	return nil
}
