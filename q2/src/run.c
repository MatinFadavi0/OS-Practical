#include <ctype.h>
#include <errno.h>
#include <sched.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/resource.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>

#include "runner.h"

#define DEFAULT_HOSTNAME "langar"

/* Runs in the child, before exec, only when --isolate is given.
 * The child is already PID 1 of the new PID namespace (created by the parent).
 * Here it moves itself into a fresh UTS namespace and only THEN sets the
 * hostname, so the host's hostname can never be modified.
 * Returns 0 on success, -1 on failure (error already printed). */
static int setup_isolated_child(const struct runner_options *options)
{
    if (unshare(CLONE_NEWUTS) < 0)
    {
        fprintf(stderr, "langar: unshare(CLONE_NEWUTS) failed: %s\n", strerror(errno));
        return -1;
    }

    const char *name = options->hostname ? options->hostname : DEFAULT_HOSTNAME;
    if (sethostname(name, strlen(name)) < 0)
    {
        fprintf(stderr, "langar: sethostname failed: %s\n", strerror(errno));
        return -1;
    }
    return 0;
}

/* Parses a CPU LIST such as "0" or "0,2" into set.
 * Only comma-separated decimal numbers are accepted: no ranges, no empty
 * items (so "1," and ",1" are rejected), no signs or spaces, and every id
 * must be below CPU_SETSIZE. Duplicates are allowed.
 * Returns 0 on success, -1 on invalid input (error already printed). */
static int parse_cpu_list(const char *list, cpu_set_t *set)
{
    CPU_ZERO(set);
    const char *p = list;
    for (;;)
    {
        if (!isdigit((unsigned char)*p))
        {
            fprintf(stderr, "langar: invalid CPU list: %s\n", list);
            return -1;
        }

        char *end;
        errno = 0;
        unsigned long cpu = strtoul(p, &end, 10);
        if (errno || cpu >= CPU_SETSIZE)
        {
            fprintf(stderr, "langar: CPU id out of range in list: %s\n", list);
            return -1;
        }
        CPU_SET(cpu, set);

        if (*end == '\0')
        {
            return 0;
        }
        if (*end != ',')
        {
            fprintf(stderr, "langar: invalid CPU list: %s\n", list);
            return -1;
        }
        p = end + 1;
    }
}

/* Runs in the child, before exec. Both settings survive execvp and are
 * inherited by the command's own children; the launcher (parent) keeps its
 * original affinity and limits. Options that were not given change nothing.
 * Returns 0 on success, -1 on failure (error already printed). */
static int apply_resource_limits(const struct runner_options *options)
{
    if (options->cpus)
    {
        cpu_set_t set;
        if (parse_cpu_list(options->cpus, &set) < 0)
        {
            return -1;
        }
        /* pid 0 = this process (the child). Fails with EINVAL when none of
         * the requested CPUs is allowed for us. */
        if (sched_setaffinity(0, sizeof(set), &set) < 0)
        {
            fprintf(stderr, "langar: sched_setaffinity failed: %s\n", strerror(errno));
            return -1;
        }
    }

    if (options->as_mib)
    {
        /* main.c already guarantees 1 <= as_mib <= 1048576, so no overflow. */
        rlim_t bytes = (rlim_t)options->as_mib * 1048576;
        struct rlimit limit = {.rlim_cur = bytes, .rlim_max = bytes};
        if (setrlimit(RLIMIT_AS, &limit) < 0)
        {
            fprintf(stderr, "langar: setrlimit(RLIMIT_AS) failed: %s\n", strerror(errno));
            return -1;
        }
    }
    return 0;
}

int run_command(char *const argv[], const struct runner_options *options)
{
    int isolate = options && options->isolate;

    /* unshare(CLONE_NEWPID) does not move the caller itself into the new PID
     * namespace; it only makes the NEXT child created by fork() become PID 1
     * there. That is why it must happen before fork(). */
    if (isolate && unshare(CLONE_NEWPID) < 0)
    {
        fprintf(stderr, "langar: unshare(CLONE_NEWPID) failed: %s\n", strerror(errno));
        return 1;
    }

    pid_t pid = fork();
    if (pid < 0)
    {
        fprintf(stderr, "langar: fork failed: %s\n", strerror(errno));
        return 1;
    }

    if (pid == 0)
    {
        /* Child: optional namespace setup, then replace this process image.
         * stdin/stdout/stderr, environment and cwd are inherited as-is. */
        if (isolate && setup_isolated_child(options) < 0)
        {
            _exit(1);
        }

        if (options && apply_resource_limits(options) < 0)
        {
            _exit(1);
        }

        execvp(argv[0], argv);

        /* Only reached if execvp failed. */
        int err = errno;
        fprintf(stderr, "langar: %s: %s\n", argv[0], strerror(err));
        _exit((err == ENOENT || err == ENOTDIR) ? 127 : 126);
    }

    /* Parent: block until the child terminates, retrying on EINTR.
     * The parent stays in the original PID namespace and sees the child by
     * its normal (host) PID. */
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
