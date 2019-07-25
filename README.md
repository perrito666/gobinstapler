# gobinstapler

A method to bundle files and folders into a go binary.

## Rationale

I often do quick tools whose UI is a crude web fron with a few endpoints to be used locally, this proved a bit annoying in go since I was not able to bundle these htmls in a way that could satisfy me, especially in terms of proper versioning.

### Components

This tool is made of two components:

* `gobinstapler`: the product of building the main package, this tool will create a fat binary containing the passed in go binary, a tar file adding recursively the passed in paths to it and a uint64 stating the go binary length to be able to access the tar afterwards.
* `prytool` package: this is a couple of functions:
    * `ListBundledFiles`: returns a map of files to their info for all the ones contained in the current running binary, optionally it receives a path to look for a binary by hand (you could point this to a different stapled binary)
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