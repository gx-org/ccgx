cmake_minimum_required (VERSION 3.24)
project (helloworld)

set (helloworld_VERSION_MAJOR 1)
set (helloworld_VERSION_MINOR 0)

include_directories(. gxdeps/github.com/gx-org/gx)

add_executable (helloworld helloworld.cc)

target_link_libraries (helloworld)

install (TARGETS helloworld DESTINATION bin)
