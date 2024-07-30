#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int main(int argc, char *argv[]) {
    if (argc != 2) {
        if (argc < 2) {
            printf("Expected exactly one command line argument, got 0\n");
        } else {
            printf("Expected exactly one command line argument, got %d\n", argc - 1);
        }
        return 1;
    }

    char *param = argv[1];

    // Random string placeholder
    const char *secretCode = "<<RANDOM>>";

    // Print the random string and the command line parameter
    printf("Hello %s! The secret code is %s.\n", param, secretCode);

    return 0;
}