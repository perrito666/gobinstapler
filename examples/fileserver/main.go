package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/perrito666/gobinstapler/prytool"
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

func main() {
	fs, err := prytool.ListStructuredFiles("")
	if err != nil {
		fmt.Printf("Something went wrong: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c := fs
		parts := strings.Split(r.URL.Path, "/")
		if parts[len(parts)-1] == "" {
			parts = parts[:len(parts)-1]
		}
		for i, part := range parts {
			if part == "" && len(parts) != 1 {
				// first one, when / preceeded, no harm jumping it
				continue
			}
			if i != len(parts)-1 {
				if f, ok := c.Folders[part]; ok {
					c = f
					continue
				}
				break
			}
			if f, ok := c.Files[part]; ok {
				tf, err := prytool.RetrieveFile("", f.FullPath)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(fmt.Sprintf("%v", err)))
					return
				}
				defer tf.Close()

				resp, err := ioutil.ReadAll(tf)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(fmt.Sprintf("%v", err)))
					return
				}
				w.Write(resp)
				return
			}
			if f, ok := c.Folders[part]; ok || i == 0 {
				if i == 0 {
					f = c
				}
				w.Write([]byte("<html><head></head><body>"))
				for _, file := range f.Files {
					w.Write([]byte(fmt.Sprintf("<p><a href='%s'>%s</a> (%d)</p>", file.FullPath, file.Name, file.Size)))
				}
				for _, folder := range f.Folders {
					w.Write([]byte(fmt.Sprintf("<p><a href='%s'>%s</a>/</p>", folder.FullPath, folder.Name)))
				}
				w.Write([]byte("</body></html>"))
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
	})
	fmt.Println("go to http://localhost:7070 to try this")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:7070"), nil))
}
