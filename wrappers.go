package gravitas

/*
#include <stdio.h>
#include "gravity_macros.h"
#include "gravity_compiler.h"
#include "gravity_delegate.h"
#include "gravity_core.h"
#include "gravity_vm.h"

extern void go_log_callback(char *description, void *xdata);
extern void go_err_callback(char *description, void *xdata);

void log_trampoline(const char *message, void *xdata) {
    go_log_callback((char *)message, xdata);
}

void error_trampoline(error_type_t error_type, const char *description, error_desc_t error_desc, void *xdata) {
    go_err_callback((char *)description, xdata);
}
*/
import "C"
