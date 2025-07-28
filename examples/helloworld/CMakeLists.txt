cmake_minimum_required (VERSION 3.24)
project (helloworld)
list(APPEND CMAKE_MODULE_PATH "${CMAKE_CURRENT_SOURCE_DIR}/gxdeps")

set (helloworld_VERSION_MAJOR 1)
set (helloworld_VERSION_MINOR 0)

include_directories (. gxdeps/github.com/gx-org/gx)

add_executable (helloworld helloworld.cc)

target_link_libraries (helloworld  ${CMAKE_SOURCE_DIR}/gxdeps/carchive.a absl_status absl_statusor)

install (TARGETS helloworld DESTINATION bin)
