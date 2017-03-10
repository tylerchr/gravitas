#include <stdio.h>
#include "gravity_macros.h"
#include "gravity_compiler.h"
#include "gravity_core.h"
#include "gravity_vm.h"

#define NEW_FUNCTION(_fptr)						(gravity_function_new_internal(NULL, NULL, _fptr, 0))
#define NEW_CLOSURE_VALUE(_fptr)				((gravity_value_t){	.isa = gravity_class_closure,		\
																	.p = (gravity_object_t *)gravity_closure_new(NULL, NEW_FUNCTION(_fptr))})

// report error callback function
void report_error(error_type_t error_type, const char *message, error_desc_t error_desc, void *xdata) {
	printf("error: %s\n", message);
	exit(0);
}

// report error callback function
void log_something(const char *message, void *xdata) {
	printf("log: %s\n", message);
}

// function to be execute inside Gravity VM		
bool my_function(gravity_vm *vm, gravity_value_t *args, uint16_t nargs, uint32_t rindex) {
	// do something useful here
	return true;
}

int main(int args, char** argv) {
	
	// Configure VM delegate		
	gravity_delegate_t delegate = {
		.log_callback = log_something,
		.error_callback = report_error,
	};
	
	// Create a new VM
	gravity_vm *vm = gravity_vm_new(&delegate);






	size_t size = 0;
	const char *input_file = "program.gravity";
	const char *source_code = file_read(input_file, &size);
	if (!source_code) {
		printf("Error loading file %s", input_file);
		exit(1);
	}

	printf("here 1\n");

	// for (int i=0; i<size; i++) {
	// 	printf("%d: %c\n", i, source_code[i]);
	// }

	// printf("%s is %d bytes long\n", input_file, size);
	
	// create compiler
	gravity_compiler_t *compiler = gravity_compiler_create(&delegate);

	printf("here 1.5\n");
	
	// compile source code into a closure
	gravity_closure_t *main_closure = gravity_compiler_run(compiler, source_code, size, 0, false);

	printf("here 2\n");

	// op is OP_COMPILE_RUN so transfer memory from compiler to VM
	gravity_compiler_transfer(compiler, vm);




	
	// Create a new class
	gravity_class_t *c = gravity_class_new_pair(vm, "MyClass", NULL, 0, 0);
	
	// Allocate and bind closures to the newly created class
	gravity_value_t closureVal = NEW_CLOSURE_VALUE(my_function);
	gravity_closure_t *closure = gravity_closure_new(vm, VALUE_AS_FUNCTION(closureVal));
	gravity_class_bind(c, "myfunc", VALUE_FROM_OBJECT(closure));
	
	// Register class inside VM
	gravity_vm_setvalue(vm, "MyClass", VALUE_FROM_OBJECT(c));





	// run the program
	printf("here\n");
	if (gravity_vm_run(vm, main_closure)) {

		gravity_value_t result = gravity_vm_result(vm);
		
		char buffer[512];
		gravity_value_dump(result, buffer, sizeof(buffer));
		printf("RESULT: %s\n", buffer);

	}

	return 0;
}