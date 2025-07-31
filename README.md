# CCGX: Generate C++ bindings from GX code

## Installation

1. Install a recent version of Go (see [All Releases](https://go.dev/dl/))
2. Install `ccgx` using the following command:
  ```
  $ GOBIN=<target folder> go install github.com/gx-org/ccgx@latest
  ```
   (if `GOBIN` is not specified, the default is `~/go/bin/ccgx`)

## Running `helloworld`

1. Create the project folder: 
    $ mkdir helloworld
    $ cd helloworld
2. Create a minimal `helloworld.gx` file such as:
    package helloworld
    
    import _ "github.com/gx-org/xlapjrt/gx"
    
    // Hello returns a constant array of two axes of size 2 and 3.
    func Hello() [2][3]float32 {
    	return [2][3]float32{
    		{1, 2, 3},
    		{4, 5, 6},
    	}
    }
3. Run 
    $ ccgx init 
   at the top of the project to create `go.mod`. `go.mod` manages all the dependencies of the project and their versions for reproducable builds.

## Disclaimer

This is not an official Google DeepMind product (experimental or otherwise), it is
just code that happens to be owned by Google-DeepMind. GX is experimental and work
in progress. As of today, we do not consider any part of the language as stable. Breaking
changes will happen on a regular basis.

You are welcome to send PR or to report bugs. We will do our best to answer but there
is no guarantee that you will get a response.
