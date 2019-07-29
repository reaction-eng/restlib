package file

import "io"

type Storage interface {
	GetArbitraryFile(id string) (io.ReadCloser, error)
	PostArbitraryFile(fileName string, parent string, file io.Reader, mime string) (string, error)
}
