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
	"io"
	"os"

	"github.com/juju/gnuflag"
	"github.com/pkg/errors"
)

var (
	// ErrNothingToStaple sould be returned when we don't have arguments for the files
	// or folders to staple to the binary.
	ErrNothingToStaple = fmt.Errorf("there is nothing to staple to the binary")
	// ErrNoBinary should be returned when the binary passed to staple things to is
	// either not a binary or not present.
	ErrNoBinary = fmt.Errorf("the first argument is not a binary")
)

type configFields struct {
	goBinary      string
	target        string
	stapleTargets []string
}

var config = configFields{}

func init() {
	gnuflag.CommandLine.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s [flags] <go-binary> <target> [fileOrFolder1 ...fileOrFolderN]:\n\n", os.Args[0])
		gnuflag.PrintDefaults()
	}

	var err error
	defer func() {
		if err != nil {
			fmt.Println(fmt.Sprintf("Failures encountered: %v\n\n", err))
			gnuflag.CommandLine.Usage()
			os.Exit(1)
		}
	}()
	err = gnuflag.CommandLine.Parse(true, os.Args[1:])
	if err != nil {
		err = errors.Wrap(err, "parsing comand arguments")
		return
	}

	args := gnuflag.CommandLine.Args()
	if len(args) < 3 {
		err = ErrNothingToStaple
		if len(args) < 1 {
			err = errors.Wrapf(err, "missing some necessary arguments: %v", ErrNoBinary)
		}
		return
	}

	config.goBinary = args[0]
	var fInfo os.FileInfo
	fInfo, err = os.Stat(args[0])
	if err != nil {
		err = errors.Wrap(err, "performing stat in the binary")
		return
	}
	if fInfo.IsDir() {
		err = errors.Wrapf(ErrNoBinary, "the passed go binary (%s) is a directory", config.goBinary)
		return
	}

	config.target = args[1]
	fInfo, tErr := os.Stat(args[1])
	if tErr == nil {
		err = errors.Wrap(err, "the target file exists, for safety reasons we don't overwrite")
		return
	}

	// TODO: check this is a binary... for some definition of elf.
	config.stapleTargets = make([]string, len(args)-2, len(args)-2)
	for i := 2; i < len(args); i++ {
		fInfo, err = os.Stat(args[i])
		if err != nil {
			err = errors.Wrap(err, "performing stat in one of the staple files")
			return
		}

		config.stapleTargets[i-2] = args[i]
	}

}

func fattenBinary() error {
	target, err := os.Create(config.target)
	if err != nil {
		return errors.Errorf("cannot create the target file: %v", err)
	}
	defer target.Close()
	goBinary, err := os.Open(config.goBinary)
	if err != nil {
		return errors.Errorf("cannot read source binary file: %v", err)
	}
	defer goBinary.Close()
	binarySize, err := io.Copy(target, goBinary)
	if err != nil {
		return errors.Errorf("copying the go binary into the target: %v", err)
	}

	_, err = buildTar(config.stapleTargets, target)
	if err != nil {
		return errors.Errorf("writing files into stapled go: %v", err)
	}
	fmt.Fprintf(target, "%020d", uint64(binarySize))
	return nil
}

func main() {
	if err := fattenBinary(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
