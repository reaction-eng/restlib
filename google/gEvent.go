// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package google

import (
	"encoding/json"
)

//Hold the struct need to create a tree
type gEvent struct {
	gFile

	//Keep a boolean if it is a file
	InfoId string `json:"infoId"`

	//Keep a boolean if it is a file
	SignupId string `json:"signupId"`

	//Keep the parent Id, unless this is root and then it is null
	ParentId string `json:"parentid"`
}

func (dir gEvent) MarshalJSON() ([]byte, error) {
	//Store a ddir copy
	type fakeItem gEvent

	//define the item type
	type ItemWithType struct {
		fakeItem
		InternalItemType string
	}

	return json.Marshal(ItemWithType{
		(fakeItem)(dir),
		"gEvent",
	})
}
