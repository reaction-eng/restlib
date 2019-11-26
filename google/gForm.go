// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package google

import (
	"encoding/json"
)

//Store the forms metadata
type FormMetaData struct {
	Title string `json:"title"`

	EmailTo []string `json:"emailTo"`

	EmailTemplate string `json:"emailTemplate"`

	EmailSubjectField string `json:"emailSubjectField"`

	DriveInfo []FormDriveInfo `json:"driveInfo"`

	RequiredPerm []string `json:"requiredPerm"`
}

//Store the forms metadata
type FormDriveInfo struct {
	SheetId   string `json:"sheetId"`
	SheetName string `json:"sheetName"`
}

//Hold the struct need to create a tree
type gForm struct {
	gFile
	//Keep the type for directory anyways
	Type string `json:"type"`

	//Keep the parent Id, unless this is root and then it is null
	ParentId string `json:"parentid"`

	//Hold the item
	Metadata FormMetaData `json:"metadata"`

	//Hold a list of Items
	JSONSchema map[string]interface{} `json:"JSONSchema"`

	//Keep the parent Id, unless this is root and then it is null
	UISchema map[string]interface{} `json:"UISchema"`
}

func (dir gForm) MarshalJSON() ([]byte, error) {
	//Store a ddir copy
	type fakeItem gForm

	//define the item type
	type ItemWithType struct {
		fakeItem
		InternalItemType string
	}

	return json.Marshal(ItemWithType{
		(fakeItem)(dir),
		"gForm",
	})
}
