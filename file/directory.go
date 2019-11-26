package file

//go:generate mockgen -destination=../mocks/mock_directory.go -package=mocks github.com/reaction-eng/restlib/file Directory

type Directory interface {
	Item

	GetType() string

	GetItems() []Item

	GetParentId() string

	ForEach(doOnFolder bool, do ItemFunc)
}

//Define a function
type ItemFunc func(item Item)
