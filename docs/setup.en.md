# Linux setup

Use Linux on your machine or in a virtual machine. Compile for that machine's
architecture; amd64 and arm64 are both supported. Native macOS is not the grading
environment: Q3 needs Linux `/proc`. Docker is not required for either problem.

On Debian/Ubuntu, install the tools once:

```sh
sudo apt update
sudo apt install build-essential python3 bash procps git
```

After setup, all assignment builds, tests, and Q3 commands run without sudo or
network access. Python 3.10 or later is required for the supplied tools; they
use only its standard library. No pip dependencies, Docker images, kernel build,
strace, or special helper program are needed.

```sh
git clone https://github.com/MohammadAminKoohi/OS-Practical.git
cd OS-Practical
git switch --detach hw02-v1.0.0
make -C assignments/hw02/q2 CC=gcc
python3 tests/hw02/test_langar.py
```

The starter compiles but is intentionally incomplete: only the CLI test should
pass, for 5/100 functional points. Implement `src/run.c` before expecting the
other tests to pass. See the [test guide](../tests/hw02/README.en.md).

For Q3, open a Bash terminal (`bash` if your shell is different) and follow the
[course handout](../assignments/hw02/handouts/Q3.fa.tex). Use the same shell for
all steps so `$!`, `$$`, and `wait` refer to the correct processes.

To keep your work on a branch after checking out the release:

```sh
git switch -c my-hw02
```

For later assignments, keep your own implementation and follow that release's
migration instructions. Do not overwrite your work with a fresh starter.
