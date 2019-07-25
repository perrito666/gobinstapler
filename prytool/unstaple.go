package prytool

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"
)

/*
MIT License

Copyright (c) 2019 Horacio Duran <horacio.duran@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/

// FileInfo contains information about a stapled file.
type FileInfo struct {
	FileInfo os.FileInfo
	FileName string
}

// SizeStringLen represents the expected size of the string containing the binary lenght
const SizeStringLen = 20

func goToFileBoundary(f *os.File) error {
	stat, err := f.Stat()
	if err != nil {
		return errors.Wrap(err, "getting file info")
	}
	sizeString := make([]byte, SizeStringLen, SizeStringLen)
	_, err = f.ReadAt(sizeString, stat.Size()-SizeStringLen)
	if err != nil {
		return errors.Wrap(err, "reading the binary size string")
	}
	size, err := strconv.ParseUint(string(sizeString), 10, 64)
	if err != nil {
		return errors.Wrap(err, "parsing size string")
	}
	_, err = f.Seek(int64(size), 0)
	return err
}

func openBinFile(binLocation string) (*os.File, error) {
	if binLocation == "" {
		var err error
		binLocation, err = os.Executable()
		if err != nil {
			return nil, errors.Wrap(err, "determining the binary path")
		}
	}
	f, err := os.Open(binLocation)
	if err != nil {
		return nil, errors.Wrap(err, "opening the binary to list files")
	}
	return f, nil
}

// ListBundledFiles returns a map of file path to File information for the files
// contained in the current stapled binary (or fails)
func ListBundledFiles(binLocation string) (map[string]FileInfo, error) {
	f, err := openBinFile(binLocation)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	err = goToFileBoundary(f)
	if err != nil {
		return nil, err
	}
	tr := tar.NewReader(f)
	result := map[string]FileInfo{}
	for {
		nextHDR, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, "reading tar headers from stapled file")
		}
		_, fName := filepath.Split(nextHDR.Name)
		result[nextHDR.Name] = FileInfo{
			FileInfo: nextHDR.FileInfo(),
			FileName: fName,
		}
	}

	return result, nil
}

// TarFile wraps a *tar.Reader and the underlying file so it can be returned
// as a file handler after seeking the right file in the tar and also provides
// the ability to close the underlying file.
type TarFile struct {
	tarFileH *tar.Reader
	fileH    *os.File
}

// Read implements io.Reader and allows to fetch one file from a tar as if it was
// a regular file
func (t *TarFile) Read(p []byte) (n int, err error) {
	return t.tarFileH.Read(p)
}

// Close implements Closer and closes the underlying file to the tar reader.
func (t *TarFile) Close() error {
	return t.fileH.Close()
}

// RetrieveFile will return a read closer for the required file or fail.
func RetrieveFile(binLocation string, filePath string) (io.ReadCloser, error) {
	f, err := openBinFile(binLocation)
	if err != nil {
		return nil, err
	}
	err = goToFileBoundary(f)
	if err != nil {
		return nil, err
	}

	tr := tar.NewReader(f)
	for {
		nextHDR, err := tr.Next()
		if err == io.EOF {
			return nil, errors.Errorf("file %s not found", filePath)
		}
		if err != nil {
			return nil, errors.Wrap(err, "reading tar files from stapled file")
		}
		if nextHDR.Name == filePath {
			return &TarFile{
				tarFileH: tr,
				fileH:    f,
			}, nil
		}
	}

}
