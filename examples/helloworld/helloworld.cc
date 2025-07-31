#include <absl/status/statusor.h>
#include <iostream>

#include "gxdeps/github.com/gx-org/ccgx/examples/helloworld/helloworld.h"
#include "gxdeps/github.com/gx-org/gx/golang/binder/ccgx/cppgx.h"
#include "gxdeps/github.com/gx-org/gx/golang/binder/cgx/cgx.h"
#include "gxdeps/github.com/gx-org/xlapjrt/cgx/cgx.cgo.h"

using github_com::gx_org::ccgx::examples::helloworld::Helloworld;
using gxlang::cppgx::Runtime;
using gxlang::cppgx::ToErrorStatus;
using std::begin;

extern "C" {
void InitModuleHelloworld();
}

absl::StatusOr<Runtime> NewRuntime() {
  InitModuleHelloworld();
  const auto bld(cgx_builder_new_static_xlapjrt());
  const auto result(cgx_runtime_new_xlapjrt(bld, "cpu"));
  if (result.error != cgx_error{}) {
    return ToErrorStatus(result.error);
  }
  return Runtime(result.runtime);
}

int main() {
  auto runtime(NewRuntime());
  if (!runtime.ok()) {
    std::cerr << "cannot create runtime: " << runtime.status() << std::endl;
    return 1;
  }
  auto device(runtime->GetDevice(0));
  if (!device.ok()) {
    std::cerr << "cannot create device: " << device.status() << std::endl;
    return 1;
  }
  auto package(Helloworld::BuildFor(device.value()));
  if (!package.ok()) {
    std::cerr << "cannot compile package: " << package.status() << std::endl;
    return 1;
  }
  auto hello(package->Hello());
  if (!hello.ok()) {
    std::cerr << hello.status() << std::endl;
  }
  auto data(hello->Acquire());
  if (!data.ok()) {
    std::cerr << data.status() << std::endl;
  }
  std::cout << "Hello: [ ";
  for (float value : data.value()) {
    std::cout << value << " ";
  }
  std::cout << "]" << std::endl;
  return 0;
}
