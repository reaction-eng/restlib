// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package google

import (
	"encoding/json"
)

type gDocument struct {
	gFile

	Type string `json:"type"`

	//Keep the Preview
	Preview string `json:"preview"`

	//Thumbnail Image
	ThumbnailUrl string `json:"thumbnail"`

	//Also Keep the parent Id
	ParentId string `json:"parentid"`
}

func (doc *gDocument) GetType() string {
	return doc.Type
}

func (doc *gDocument) GetParentId() string {
	return doc.ParentId
}

func (doc *gDocument) GetPreview() string {
	return doc.Preview
}

func (doc *gDocument) GetThumbnailUrl() string {
	return doc.ThumbnailUrl
}

func (doc gDocument) MarshalJSON() ([]byte, error) {
	//Store a dir copy
	type fakeItem gDocument

	//define the item type
	type ItemWithType struct {
		fakeItem
		InternalItemType string
	}

	return json.Marshal(ItemWithType{
		(fakeItem)(doc),
		"gDocument",
	})
}
