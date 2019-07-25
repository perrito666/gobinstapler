package main

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
	return addToTar(filesToInclude, tWriter)
}

// buildTar will try to write the files passed into the passed ioWriter in tar format
// if the writing fails some data might already be writte.
func addToTar(filesToInclude []string, tWriter *tar.Writer) (int64, error) {
	var contentsSize int64
	tWriter.Flush()
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

	return contentsSize, tWriter.Flush()
}
