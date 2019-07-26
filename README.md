# gobinstapler

A tool to make fat binaries by "stapling" files to it using tar.

This produces a version of your go executable that contains files and folders in the binary but they are added post compilation and can be changed without re-compiling the go binary. Functions are provided so you can create code that takes advantage of these binaries as if they were on a filesystem path (you get `io.Reader` handlers).

The input file is copied into a target file along with one or more files or folders passed as arguments for a result binary shaped like this:

```
 ____________________________target executable_________________________
/                                                                       \
[ gobinary | tar formatted files and folders | footer Table of contents ]
```

Additionally you can pass a `--relative` flag to make the target filesystem relative to that folder (meaning it gets stripped)

You can see an example in this same repo under `examples/fileserver`:
* go install your gobinstapler (or make it available in the path)
* invoke `run.sh`  in the fileserver example folder passing as parameter at least one file or folder for the resulting server and then execute the resulting `selfserve` and go to your browser's `http://localhost:7070`


## Rationale

I often do quick tools whose UI is a crude web front-end with a few endpoints to be used locally, this proved a bit annoying in go since I was not able to bundle these htmls in a way that could satisfy me, especially in terms of proper versioning.
There most likely are tools that do this or similar things, the one that comes to mind is [statik](https://github.com/rakyll/statik) which I have used before but I wanted three things:
* I fancied writting the thing (and that by itself should be enough)
* I wanted the tool to only affect the final binary
* I wanted to have as little stuff in mem for this as possible so I can use large files, this way you get a file handler to a section of your tar I think so you get a sort of mmap

### Components

This tool is made of two components:

* `gobinstapler`: the product of building the main package, this tool will create a fat binary containing the passed in go binary, a tar file adding recursively the passed in paths to it and a uint64 stating the go binary length to be able to access the tar afterwards.
* `prytool` package: this is three functions:
    * `ListBundledFiles`: returns a map of files to their info for all the ones contained in the current running binary, optionally it receives a path to look for a binary by hand (you could point this to a different stapled binary)
    * `ListStructuredFiles`: returns a `prytool.Folder` object containing a tree of `prytool.File` and `prytool.Folder` for you to use in a more natural way.
    * `RetrieveFile`: returns a file from the stapled binary by path (the path being the key in the map returned by `ListBundledFiles`)

### Usage 

```bash
gobinstapler gobinary targetbinary folder1 folder2 foldern file
```

where:

* `gobinstapler` is the command built in the main package
* `gobinary` is the path to the go binary we want to staple the files or folders to
* `targetbinary` where the final stapled file will be stored (this does not modify the original binary just in case)
* `folder1...folderN.file1..fileN` a mixed list of files and folders to be added, the path of these, as passed, will be used (ie if we pass something in `/home/user/foler` the path will be that but if we are in `/home/user` and just pass `folder` that will be it, they are relative). Passing a folder and then a file in that folder is undefined content (I have not tested)

### Testing

Inside the testing folder, you can find a makefile, running `make clean; make test` will build `gobinstapler` then a test binary, staple some test folders and run the binary that is a suite of tests itself.
