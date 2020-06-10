package file

import (
	"time"
)

type Listing struct {
	//Keep a boolean if it is a file
	Id string `json:"id"`

	//Hold the item
	Name string `json:"name"`

	//Hold if we should hide the item
	HideListing bool `json:"hideListing"`

	//Keep a date if useful
	Date *time.Time `json:"date"`

	Listings []Listing `json:"listings"`
	Items    []Item    `json:"items"`

	//Keep the parent Id, unless this is root and then it is null
	ParentId string `json:"parentId"`
}

func NewListing() *Listing {
	return &Listing{
		Listings: make([]Listing, 0),
		Items:    make([]Item, 0),
	}
}

//Define a function
type ItemFunc func(item *Item)

func (dir *Listing) ForEach(doOnFolder bool, do ItemFunc) {
	if dir == nil {
		return
	}

	//March over each item in the dir
	for _, item := range dir.Items {
		do(&item)
	}

	for _, subDir := range dir.Listings {
		subDir.ForEach(doOnFolder, do)
	}
}
