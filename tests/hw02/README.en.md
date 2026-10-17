# HW02 public tests

[فارسی](README.fa.md)

From the repository root, build and run:

```sh
make -C assignments/hw02/q2 clean
make -C assignments/hw02/q2 CC=gcc
python3 tests/hw02/test_langar.py
```

The untouched starter should build and score **5/100**, passing only `cli_usage`.
The other failures are the student's implementation task.

## Functional points (100 total within Q2)

- `cli_usage`: 5 — missing command/separator and invalid subcommand.
- `absolute_execution`: 5 — exact command output, no launcher chatter.
- `path_lookup`: 5 — command lookup using a temporary PATH.
- `argument_boundaries`: 10 — spaces, empty values, UTF-8 and literal punctuation.
- `stdin_inheritance`: 5 — command receives its input unchanged.
- `stdout_stderr`: 5 — streams stay separate.
- `working_directory`: 5 — caller's directory is preserved.
- `environment`: 5 — environment is inherited.
- `direct_child`: 15 — probe PID differs from launcher; probe PPID is launcher PID.
- `parent_waits`: 10 — launcher stays alive while child waits for input.
- `normal_exit_status`: 10 — normal statuses 0, 7, 42, and 255.
- `signal_exit_status`: 5 — SIGTERM maps to 143.
- `command_not_found`: 5 — missing program/path maps to 127 and stderr.
- `permission_denied`: 5 — non-executable file maps to 126 and stderr.
- `not_a_directory`: 5 — ENOTDIR maps to 127 and stderr.

Each case earns its listed points when all its checks pass. This is the functional
score, not the entire course grade. Code/API requirements and explanations are
reviewed separately; videos and reports are required. Resource-failure handling
and actual use of waitpid/reaping are reviewed in source. The behavioral waiting
test is not proof of a particular API call. Final evaluation can vary inputs
while preserving these specified capabilities and weights.

## Logs and machine-readable results

Save the process's real exit status while capturing both output streams:

```sh
python3 tests/hw02/test_langar.py > assignments/hw02/q2/test-results.txt 2>&1
status=$?
cat assignments/hw02/q2/test-results.txt
echo "$status"
```

Exit code 0 means all tests passed; 1 means functional failures; 2 means test setup
failed (for example, executable/compiler missing). A failed run still provides a
valid result log to submit. Do not fabricate passing results.

To test an extracted submission or save JSON:

```sh
python3 tests/hw02/test_langar.py --binary /path/to/submission/q2/langar --json /tmp/hw02-result.json
```

JSON contains `score`, `total`, and per-case `id`, `points`, `earned`, `passed`,
and `error`. This lets course tooling aggregate many submissions without parsing
human-readable logs. Tests are independent of the submitted Makefile/test files.
`--cc gcc` selects the compiler used to build the test probe.

The probe source is a workload, not a runner solution. It is compiled in a
temporary directory for the current architecture. Each test execution has a
5-second watchdog; the release is checked on Linux. The readiness handshake
keeps the waiting test independent of a guessed sleep duration. Test-owned
process groups are cleaned up after success/failure.
