package file

//go:generate mockgen -destination=../mocks/mock_document.go -package=mocks github.com/reaction-eng/restlib/file Document

type Document interface {
	Item

	Type() string

	Preview() string

	ThumbnailUrl() string

	ParentId() string
}
