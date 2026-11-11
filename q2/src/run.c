#include <errno.h>
#include <stdio.h>
#include <string.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>

#include "runner.h"

/* Retain your own previous runner and extend its setup/wait paths.
 * The handout defines this version's required behavior. This is not a solution.
 * HW3 migration: add the options parameter declared in runner.h.
 */
int run_command(char *const argv[], const struct runner_options *options)
{
    /* TODO (HW03): use options (isolate, hostname, cpus, rootfs, image, as_mib). */
    (void)options;

    pid_t pid = fork();
    if (pid < 0)
    {
        fprintf(stderr, "langar: fork failed: %s\n", strerror(errno));
        return 1;
    }

    if (pid == 0)
    {
        /* Child: replace this process image with the requested program.
         * stdin/stdout/stderr, environment and cwd are inherited as-is. */
        execvp(argv[0], argv);

        /* Only reached if execvp failed. */
        int err = errno;
        fprintf(stderr, "langar: %s: %s\n", argv[0], strerror(err));
        _exit((err == ENOENT || err == ENOTDIR) ? 127 : 126);
    }

    /* Parent: block until the child terminates, retrying on EINTR. */
    int status;
    pid_t result;
    do
    {
        result = waitpid(pid, &status, 0);
    } while (result < 0 && errno == EINTR);

    if (result < 0)
    {
        fprintf(stderr, "langar: waitpid failed: %s\n", strerror(errno));
        return 1;
    }

    if (WIFEXITED(status))
    {
        return WEXITSTATUS(status);
    }
    if (WIFSIGNALED(status))
    {
        return 128 + WTERMSIG(status);
    }
    return 1;
}
