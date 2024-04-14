package template

import (
	"archive/tar"
	"embed"
	"fmt"
	"io"
)

const (
	Replacement = "github.com/go-chocolate/example"
)

//go:embed assets/template.tar
var fs embed.FS

type File struct {
	Filename string
	IsDir    bool
	Content  []byte
}

func Read() ([]*File, error) {
	temp, err := fs.Open("assets/template.tar")
	if err != nil {
		return nil, err
	}
	defer temp.Close()

	re := tar.NewReader(temp)

	var files []*File

	for {
		head, err := re.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		file := &File{Filename: head.Name}
		switch head.Typeflag {
		case tar.TypeDir:
			file.IsDir = true
		case tar.TypeReg:
			file.Content = make([]byte, head.Size)
			if _, err := re.Read(file.Content); err != nil && err != io.EOF {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unsupported tar flag: %v", head.Typeflag)
		}
		files = append(files, file)
	}
	return files, nil
}
