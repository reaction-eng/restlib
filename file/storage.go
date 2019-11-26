package file

//go:generate mockgen -destination=../mocks/mock_item.go -package=mocks github.com/reaction-eng/restlib/file Storage

import "io"

type Storage interface {
	GetArbitraryFile(id string) (io.ReadCloser, error)
	PostArbitraryFile(fileName string, parent string, file io.Reader, mime string) (string, error)

	BuildFileHierarchy(dirId string, buildPreview bool, includeFilter func(fileType string) bool) Directory
	BuildFormHierarchy(dirId string) Directory
	GetFilePreview(id string) string
	GetFileThumbnailUrl(id string)
	GetFileHtml(id string) string
	GetMostRecentFileInDir(dirId string) (io.ReadCloser, error)
	GetFileAsInterface(id string, inter interface{})
	GetFirstFileMatching(dirId string, name string) (io.ReadCloser, error)
}
