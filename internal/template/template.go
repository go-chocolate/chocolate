package template

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"embed"
	"fmt"
	"io"
	"os"
	"path"
)

const (
	replacement = "github.com/go-chocolate/example"
)

//go:embed assets/template.tar.gz
var fs embed.FS

type File struct {
	Filename string
	IsDir    bool
	Content  []byte
}

func Read() ([]*File, error) {
	temp, err := fs.Open("assets/template.tar.gz")
	if err != nil {
		return nil, err
	}
	defer temp.Close()

	gz, err := gzip.NewReader(temp)
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	re := tar.NewReader(gz)

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

type Option struct {
	Output string
	Module string
}

func Write(file *File, opt Option) error {
	filename := path.Clean(fmt.Sprintf("%s/%s", opt.Output, file.Filename))
	if file.IsDir {
		return os.Mkdir(filename, 0777)
	}
	var content = file.Content
	if opt.Module != "" {
		content = bytes.ReplaceAll(content, []byte(replacement), []byte(opt.Module))
	}
	return os.WriteFile(filename, content, 0644)
}
