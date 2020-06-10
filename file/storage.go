package file

//go:generate mockgen -destination=../mocks/mock_storage.go -package=mocks github.com/reaction-eng/restlib/file Storage

import "io"

type Storage interface {
	GetArbitraryFile(id string) (io.ReadCloser, error)
	PostArbitraryFile(fileName string, parent string, file io.Reader, mime string) (string, error)

	BuildListing(dirId string, previewLength int, includeFilter func(fileType string) bool) (*Listing, error)
	GetFilePreview(id string, previewLength int) string
	GetFileThumbnailUrl(id string) string
	GetFileHtml(id string) string
	GetMostRecentFileInDir(dirId string) (io.ReadCloser, error)
	GetFileAsInterface(id string, inter interface{}) error
	GetFirstFileMatching(dirId string, name string) (io.ReadCloser, error)
}
