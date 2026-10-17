#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

/* A test workload, not a launcher implementation. Compile on the test machine. */
int main(int argc, char *argv[])
{
    if (argc < 2) {
        return 64;
    }
    if (strcmp(argv[1], "args") == 0) {
        for (int i = 2; i < argc; ++i) {
            printf("%zu:%s\n", strlen(argv[i]), argv[i]);
        }
    } else if (strcmp(argv[1], "identity") == 0) {
        printf("%ld %ld\n", (long)getpid(), (long)getppid());
    } else if (strcmp(argv[1], "wait") == 0) {
        puts("READY");
        fflush(stdout);
        return getchar() == 'x' ? 0 : 65;
    } else if (strcmp(argv[1], "exit") == 0 && argc == 3) {
        return atoi(argv[2]);
    } else if (strcmp(argv[1], "signal") == 0) {
        return raise(SIGTERM) == 0 ? 66 : 67;
    } else if (strcmp(argv[1], "streams") == 0) {
        fputs("out\n", stdout);
        fputs("err\n", stderr);
    } else if (strcmp(argv[1], "stdin") == 0) {
        int c;
        while ((c = getchar()) != EOF) {
            if (putchar(c) == EOF) {
                return 68;
            }
        }
        return ferror(stdin) ? 69 : 0;
    } else if (strcmp(argv[1], "cwd") == 0) {
        char *cwd = getcwd(NULL, 0);
        if (cwd == NULL) {
            return 70;
        }
        puts(cwd);
        free(cwd);
    } else if (strcmp(argv[1], "env") == 0) {
        const char *value = getenv("LANGAR_TEST_VALUE");
        if (value == NULL) {
            return 71;
        }
        puts(value);
    } else {
        return 64;
    }
    return 0;
}
