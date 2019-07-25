package main

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

import (
	"archive/tar"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func buildTar(filesToInclude []string, tarFile io.Writer) (int64, error) {
	tWriter := tar.NewWriter(tarFile)
	defer tWriter.Close()
	return addToTar(filesToInclude, tWriter)
}

// buildTar will try to write the files passed into the passed ioWriter in tar format
// if the writing fails some data might already be writte.
func addToTar(filesToInclude []string, tWriter *tar.Writer) (int64, error) {
	var contentsSize int64
	for _, filePath := range filesToInclude {
		fInfo, err := os.Stat(filePath)
		if err != nil {
			return 0, errors.Wrap(err, "accessing file to include in stapled file")
		}
		// Header
		// for now we just dereference symlinks
		hdr, err := tar.FileInfoHeader(fInfo, "")
		if err != nil {
			return 0, errors.Wrap(err, "creating header for file")
		}
		hdr.Name = filePath
		err = tWriter.WriteHeader(hdr)
		if err != nil {
			return 0, errors.Wrapf(err, "writing header information for %s", filePath)
		}
		// File
		if !fInfo.IsDir() {
			contents, err := ioutil.ReadFile(filePath)
			if err != nil {
				return 0, errors.Wrap(err, "reading file to add into tar")
			}
			writen, err := tWriter.Write(contents)
			if err != nil {
				return 0, errors.Wrap(err, "writing file contents into tar")
			}

			// Fixme: this is wrong?
			contentsSize += int64(writen)
			continue
		}
		dirContents, err := filepath.Glob(filepath.Join(filePath, "*"))
		if err != nil {
			return 0, errors.Wrapf(err, "trying to read the contents of %s", filePath)
		}
		//for i := range dirContents {
		//	dirContents[i] = filepath.Join(filePath, dirContents[i])
		//}
		writen, err := addToTar(dirContents, tWriter)
		if err != nil {
			return 0, errors.Wrapf(err, "adding to tar the contents of %s", filePath)
		}
		contentsSize += writen
	}

	return contentsSize, nil
}
