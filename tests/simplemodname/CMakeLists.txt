cmake_minimum_required (VERSION 3.24)
project (simplemodname)
list(APPEND CMAKE_MODULE_PATH "${CMAKE_CURRENT_SOURCE_DIR}/gxdeps")

include_directories (. gxdeps/github.com/gx-org/gx)

include(${CMAKE_CURRENT_SOURCE_DIR}/gxdeps/github.com/gx-org/gx/golang/binder/ccgx/CMakeLists.txt)
include(${CMAKE_CURRENT_SOURCE_DIR}/gxdeps/simplemodname/CMakeLists.txt)

add_executable (simplemodname simplemodname.cc)

target_link_libraries (simplemodname  ${CMAKE_SOURCE_DIR}/gxdeps/carchive.a $ENV{GOPJRT_INSTALL_DIR}/lib/libgomlx_xlabuilder.a absl_status absl_statusor ccgx simplemodname_bindings)

install (TARGETS simplemodname DESTINATION bin)
