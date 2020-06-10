package file

import (
	"time"
)

type Item struct {
	//Keep a boolean if it is a file
	Id string `json:"id"`

	//Hold the item
	Name string `json:"name"`

	//Hold if we should hide the item
	HideListing bool `json:"hideListing"`

	//Keep a date if useful
	Date *time.Time `json:"date"`

	Type string `json:"type"`

	//Keep the Preview
	Preview string `json:"preview"`

	//Thumbnail Image
	ThumbnailUrl string `json:"thumbnail"`

	//Also Keep the parent Id
	ParentId string `json:"parentId"`
}

// ByAge implements sort.Interface for []Person based on
// the Age field.
type ByDate []Item

func (a ByDate) Len() int      { return len(a) }
func (a ByDate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool {
	//Get both dates
	iDate := a[i].Date
	jDate := a[j].Date

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
