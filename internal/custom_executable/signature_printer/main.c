#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int main(int argc, char *argv[]) {
    char *randomCode = argv[1];

    // Random string placeholder
    const char *secretCode = "<<RANDOM>>";

    printf("Program was passed %d args (including program name).\n", argc);
    printf("Arg #0 (program name): %s\n", argv[0]);
    printf("Arg #1: %s\n", randomCode);
    printf("Program Signature: %s\n", secretCode);

    return 0;
}
