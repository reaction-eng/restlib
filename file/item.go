package file

//go:generate mockgen -destination=../mocks/mock_storage.go -package=mocks github.com/reaction-eng/restlib/file Item

import (
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
