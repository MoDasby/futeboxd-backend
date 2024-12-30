package upload

import "io"

type File struct {
	Name        string
	Size        int64
	Content     io.ReadSeeker
	ContentType string
}
