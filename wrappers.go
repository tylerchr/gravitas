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
extern char* go_precode_callback(void *xdata);
extern char* go_loadfile_callback(const char *file, size_t *size, uint32_t *fileid, void *xdata);

void log_trampoline(const char *message, void *xdata) {
    go_log_callback((char *)message, xdata);
}

void error_trampoline(error_type_t error_type, const char *description, error_desc_t error_desc, void *xdata) {
    go_err_callback((char *)description, xdata);
}

const char* precode_trampoline(void *xdata) {
    return (const char*)go_precode_callback(xdata);
}

const char* loadfile_trampoline(const char *file, size_t *size, uint32_t *fileid, void *xdata) {
    return (const char*)go_loadfile_callback(file, size, fileid, xdata);;
}
*/
import "C"
