#include <absl/status/statusor.h>
#include <iostream>

#include "gxdeps/github.com/gx-org/gx/golang/binder/ccgx/cppgx.h"
#include "gxdeps/github.com/gx-org/xlapjrt/cgx/cgx.cgo.h"
#include <gxdeps/github.com/gx-org/gx/golang/binder/cgx/cgx.h>

using gxlang::cppgx::Runtime;
using gxlang::cppgx::ToErrorStatus;

absl::StatusOr<Runtime> NewRuntime() {
  const auto bld(cgx_builder_new_static_xlapjrt());
  const auto result(cgx_runtime_new_xlapjrt(bld, "cpu"));
  if (result.error != cgx_error{}) {
    return ToErrorStatus(result.error);
  }
  return Runtime(result.runtime);
}

int main() {
  auto rtm(NewRuntime());
  if (!rtm.ok()) {
    std::cout << rtm.status() << std::endl;
  }
  std::cout << "Runtime ok" << std::endl;
  return 0;
}
