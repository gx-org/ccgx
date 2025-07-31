# CCGX: Generate C++ bindings from GX code

## Installation

1. Install a recent version of Go (see [All Releases](https://go.dev/dl/))
2. Install `ccgx` using the following command:
    ```
    $ GOBIN=<target folder> go install github.com/gx-org/ccgx@latest
    ```
   (if `GOBIN` is not specified, the default is `~/go/bin/ccgx`)
3. Install [gopjrt](https://github.com/gomlx/gopjrt), the XLA backend used by GX:
    ```
    $ export GOPJRT_NOSUDO=true
    $ export GOPJRT_INSTALL_DIR=$HOME/gopjrtbin
    $ curl -sSf https://raw.githubusercontent.com/gomlx/gopjrt/main/cmd/install_linux_amd64.sh | bash
    ```
    Note that `GOPJRT_INSTALL_DIR` is going to be used later in `CMakeLists.txt`. (Check the [install_linux_amd64.sh](https://github.com/gomlx/gopjrt/blob/main/cmd/install_linux_amd64.sh)).

## Running `helloworld`

This example explains how to run the example in [ccgx/examples/helloworld](https://github.com/gx-org/ccgx/blob/main/examples/helloworld).

1. Create the project folder: 
    ```
    $ mkdir helloworld
    $ cd helloworld
    ```
2. Create a minimal `helloworld.gx` file such as:
    ```go
    package helloworld
    
    import _ "github.com/gx-org/xlapjrt/gx"
    
    // Hello returns a constant array of two axes of size 2 and 3.
    func Hello() [2][3]float32 {
    	return [2][3]float32{
    		{1, 2, 3},
    		{4, 5, 6},
    	}
    }
    ```
    Note the `import` which adds a dependency to the XLA backend. See [ccgx/examples/helloworld/helloworld.gx](https://github.com/gx-org/ccgx/blob/main/examples/helloworld/helloworld.gx) as a reference.
3. Run the following command to create `go.mod` and `go.sum`:
    ```
    $ ccgx mod init helloworld
    ```
   These files manage all the dependencies of the project and their versions for reproducable builds. After new dependencies are added or removed in GX source files, run:
    ```
    $ ccgx mod tidy 
    ```
   to update `go.mod` from the latest imports in the GX source files.
4. Run the following command to generate a corresponding C++ source and header files:
    ```
    $ ccgx bind
    ```
   The files are generated in the `gxdeps` folder.
5. Create the C++ file [helloworld.cc](https://github.com/gx-org/ccgx/blob/main/examples/helloworld/helloworld.cc) and its [CMakeLists.txt](https://github.com/gx-org/ccgx/blob/main/examples/helloworld/CMakeLists.txt)
6. Compile and run the project with `cmake`:
    ```
    $ mkdir build
    $ cd build
    $ cmake ..
    $ ./helloworld
    ```

## Disclaimer

This is not an official Google DeepMind product (experimental or otherwise), it is
just code that happens to be owned by Google-DeepMind. GX is experimental and work
in progress. As of today, we do not consider any part of the language as stable. Breaking
changes will happen on a regular basis.

You are welcome to send PR or to report bugs. We will do our best to answer but there
is no guarantee that you will get a response.
