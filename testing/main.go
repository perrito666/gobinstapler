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
	"fmt"
	"os"

	"github.com/perrito666/gobinstapler/prytool"
)

var expectedFolders = map[string]bool{
	"testfolder/subfolder": true,
	"testfolder":           true,
}
var expectedFiles = map[string]bool{
	"testfolder/subfolder/file3.txt": true,
	"testfolder/file1.txt":           true,
	"testfolder/file2.txt":           true,
}

var expectedContents = map[string]string{
	"testfolder/subfolder/file3.txt": "a nexted file also with text",
	"testfolder/file1.txt":           "A file with some text",
	"testfolder/file2.txt":           "another file with text",
}

func main() {
	fs, err := prytool.ListBundledFiles("")
	if err != nil {
		fmt.Printf("FAIL: could not find bundled files: %v\n", err)
		os.Exit(1)
	}
	if len(fs) != 5 {
		fmt.Printf("FAIL: expected 5 files, got: %d\n", len(fs))
		os.Exit(1)
	}
	for k, v := range fs {
		if v.FileInfo.IsDir() {
			if _, ok := expectedFolders[k]; !ok {
				fmt.Printf("FAIL: folder %s is not expected\n", k)
				os.Exit(1)
			}
			fmt.Printf("Found folder: %s\n", k)
			continue
		}
		if _, ok := expectedFiles[k]; !ok {
			fmt.Printf("FAIL: file %s is not expected\n", k)
			os.Exit(1)
		}
		fmt.Printf("Found file: %s\n", k)
		rc, err := prytool.RetrieveFile("", k)
		if err != nil {
			if rc != nil {
				rc.Close()
			}
			fmt.Printf("FAIL: retrieving %s file contents: %v\n", k, err)
			os.Exit(1)
		}
		c := make([]byte, v.FileInfo.Size()-1, v.FileInfo.Size()-1)
		_, err = rc.Read(c)
		if err != nil {
			if rc != nil {
				rc.Close()
			}
			fmt.Printf("FAIL: reading %s file contents: %v\n", k, err)
			os.Exit(1)
		}
		rc.Close()
		if expectedContents[k] != string(c) {
			fmt.Printf("FAIL: content missmatch, expected %q but got %q\n", expectedContents[k], string(c))
			os.Exit(1)
		}
		fmt.Printf("File contents: %q\n", string(c))
	}
}
