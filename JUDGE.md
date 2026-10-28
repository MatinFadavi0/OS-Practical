# Running the HW2 judge

This repository contains starter code and judging tools. Follow the course
handout for assignment requirements and submission rules. For installation and
testing on Linux, macOS or Windows, see [Langar Lab](LANGAR-LAB.md).

The commands below run the judge directly inside Linux. After installing the
tools, building and testing HW2 require neither sudo nor a network connection.

## Prerequisites

Use Linux with GCC, Make and Python 3.10+. The judge uses only the Python
standard library; no pip packages are needed. On Debian/Ubuntu:

```sh
sudo apt update
sudo apt install build-essential python3 git
```

## Download, build and test

```sh
git clone https://github.com/MohammadAminKoohi/OS-Practical.git
cd OS-Practical
make -C q2 clean
make -C q2 CC=gcc
python3 judge/test_langar.py
```

The executable is `q2/langar`. Implement `q2/src/run.c`; the untouched starter
builds and scores **5/100**, which is expected. Each test prints PASS or FAIL,
its points and any failure details. The total functional score is 100.

The judge compiles `probe.c` in a temporary directory; do not build it manually
or use it in your implementation. Each test has a five-second timeout.

Judge exit codes:

- `0`: all tests passed.
- `1`: one or more functional tests failed.
- `2`: setup or usage error, such as a missing executable or compiler.

## Save logs and JSON results

Run from the repository root. Capture both output streams while preserving the
judge's exit status:

```sh
python3 judge/test_langar.py > q2/test-results.txt 2>&1
status=$?
cat q2/test-results.txt
echo "$status"
```

Save per-test results or test another executable:

```sh
python3 judge/test_langar.py --json /tmp/langar-result.json
python3 judge/test_langar.py --binary /path/to/submission/q2/langar --cc gcc
```

JSON includes `score`, `total` and, for each test, `id`, `points`, `earned`,
`passed` and `error`. `--binary` may point outside this repository; the judge
still uses this repository's tests. `--cc` selects the probe compiler.

After cleaning the built executable and objects, run
`python3 judge/check_submission.py /path/to/extracted-submission` to check the
required submission files. This checks layout and basic file formats only;
follow the separately provided course handout for requirements.
