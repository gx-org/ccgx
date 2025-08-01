cmake_minimum_required (VERSION 3.24)
project (urlmodname)
list(APPEND CMAKE_MODULE_PATH "${CMAKE_CURRENT_SOURCE_DIR}/gxdeps")

include_directories (. gxdeps/github.com/gx-org/gx)

include(${CMAKE_CURRENT_SOURCE_DIR}/gxdeps/github.com/gx-org/gx/golang/binder/ccgx/CMakeLists.txt)
include(${CMAKE_CURRENT_SOURCE_DIR}/gxdeps/github.com/gx-org/ccgx/tests/urlmodname/CMakeLists.txt)

add_executable (urlmodname urlmodname.cc)

target_link_libraries (urlmodname  ${CMAKE_SOURCE_DIR}/gxdeps/carchive.a $ENV{GOPJRT_INSTALL_DIR}/lib/libgomlx_xlabuilder.a absl_status absl_statusor ccgx urlmodname_bindings)

install (TARGETS urlmodname DESTINATION bin)
