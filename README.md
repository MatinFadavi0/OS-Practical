# Operating Systems Practical Assignments

Student starter code, public judges and development tools for the Operating
Systems practical assignments. The course handouts are distributed separately
and define each assignment's requirements, deadlines and submission rules.

## Assignment structure

The practical assignments have three questions:

- **Q1:** a separate course exercise managed by the course staff.
- **Q2 — Langar:** an incremental container runtime written in C. Extend your
  own implementation each week and preserve the behavior from earlier stages.
- **Q3 — Linux:** a short investigation of Linux processes, performance,
  memory or security. Run the experiment, interpret the observations and
  explain your reasoning.

Q2 and Q3 include written answers and video explanations as specified in the
handout. Q2 has automated functional tests; Q3 observations and explanations
are evaluated manually. Keep the submitted `q2/` and `q3/` directories directly
under your submission root, following that week's required file list.

## Langar

Langar is a teaching project that builds a small Linux container runtime one
operating-system concept at a time. It begins with a command launcher:

```sh
./q2/langar run -- /bin/echo hello
```

The planned progression from HW2 to HW8 is:

- **HW2 — Processes:** create a child, execute a command, wait and return its status.
- **HW3 — Namespaces:** add PID and hostname isolation.
- **HW4 — Resource controls:** restrict CPU affinity and virtual address space.
- **HW5 — Mounts:** create a private mount namespace and a matching `/proc` view.
- **HW6 — Root filesystem:** execute commands inside a supplied root filesystem.
- **HW7 — Images:** assemble local OCI image layers and apply image configuration.
- **HW8 — Process lifecycle:** supervise workloads as PID 1, forward signals,
  reap children and clean up remaining processes and runtime state.

By the end, Langar can run supported Linux image workloads with isolation,
resource controls and lifecycle management. Each week's handout defines the
supported subset and exact behavior. Starter updates and tests are released
incrementally; use the release announced for your assignment.

## Repository layout

```text
.
├── README.md                    # Course project overview
├── LANGAR-LAB.md                # Helper installation, usage and troubleshooting
├── JUDGE.md                     # Direct judge commands for this release
├── q2/
│   ├── Makefile                 # Build the C runtime
│   ├── include/runner.h         # Supplied interface
│   └── src/                    # Parser, runner and milestone-specific helpers
├── judge/                      # Public tests, probes and submission checker
└── tools/
    ├── install-langar-lab.sh    # macOS/Linux helper installer
    ├── install-langar-lab.ps1   # Windows PowerShell helper installer
    └── langar-lab/             # Helper source, embedded judges and checksums
```

Additional Q2 helpers and Q3 workloads appear as their assignments are released.
In HW2, complete `q2/src/run.c`; `q2/src/main.c` supplies the command-line parser.
Later releases extend the interface and add files. Merge these updates into
your previous implementation while retaining your working code.

## Getting started

```sh
git clone https://github.com/MohammadAminKoohi/OS-Practical.git
cd OS-Practical
```

Install the prerequisites in the [Langar Lab guide](LANGAR-LAB.md), then on
macOS/Linux run:

```sh
sh tools/install-langar-lab.sh
./.langar-lab/bin/langar-lab setup
./.langar-lab/bin/langar-lab test hw2 --source ./q2 --json ./hw2-result.json
```

On Windows, follow the guide's [PowerShell instructions](LANGAR-LAB.md#2-install-and-run).
Replace `hw2` with the matching assignment when it is released.

Langar Lab builds and tests your code inside Linux. Linux users default to
**native** execution, with no VM download. macOS and Windows use a local Ubuntu
VM through **Multipass**. The helper copies your sources, builds them, runs the
public judge and saves results and logs. Its Go source is available under
`tools/langar-lab/`; prebuilt executables are provided through
[GitHub Releases](https://github.com/MohammadAminKoohi/OS-Practical/releases/tag/langar-lab-v0.2.0).

For building and running the judge directly inside Linux, see [JUDGE.md](JUDGE.md).

## Testing and submission

Public functional tests report case results and a score out of 100. The
unfinished HW2 starter builds and scores **5/100**; implement the runner to
pass the remaining tests. Build and environment errors are reported separately
and receive no functional score. Final grading reruns the course judge on
your submitted code; local results help you find and fix problems.

Before submitting, follow the handout's file layout, include the required
answers and videos, and remove build products. The submission checker in
`judge/` checks the required layout and basic file formats. Full judge commands,
JSON results and log instructions are documented in [JUDGE.md](JUDGE.md).
