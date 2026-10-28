"""Linux-only build/test worker embedded in the Langar Lab host binary."""

import argparse
import json
import os
from pathlib import Path
import platform
import shutil
import signal
import subprocess
import sys
import time


def start_time(pid):
    try:
        return Path(f"/proc/{pid}/stat").read_text().rsplit(")", 1)[1].split()[19]
    except (OSError, IndexError):
        return None


def terminate_group(process):
    if process is None:
        return
    def send(sig):
        try:
            os.killpg(process.pid, sig)
        except ProcessLookupError:
            pass
        except PermissionError:
            # The HW3+ judge runs via sudo; the unprivileged worker may need
            # sudo to signal precisely the process group that it started.
            result = subprocess.run(['sudo', '-n', '/bin/kill', '-'+str(int(sig)),
                                     '--', '-'+str(process.pid)], capture_output=True,
                                    timeout=5)
            if result.returncode and process.poll() is None:
                raise RuntimeError('Cannot stop privileged course judge group')
    send(signal.SIGTERM)
    try:
        process.wait(timeout=1)
    except subprocess.TimeoutExpired:
        pass
    send(signal.SIGKILL)
    process.wait()


def environment(cc):
    if sys.platform != "linux":
        raise RuntimeError("The worker must run inside Linux")
    if sys.version_info < (3, 10):
        raise RuntimeError("Python 3.10 or newer is required")
    for program in ("make", cc):
        if shutil.which(program) is None:
            raise RuntimeError(f"Missing program: {program}")
    version = subprocess.run([cc, "--version"], text=True, capture_output=True,
                             timeout=10, check=True).stdout.splitlines()[0]
    distribution = {}
    for line in Path("/etc/os-release").read_text().splitlines():
        if "=" in line:
            key, value = line.split("=", 1)
            distribution[key] = value.strip('"')
    return dict(system=platform.system(), kernel=platform.release(),
                architecture=platform.machine(), python=platform.python_version(),
                compiler=version, distribution=distribution.get("PRETTY_NAME", "Linux"))


def cancel(work):
    work.mkdir(parents=True, exist_ok=True)
    (work / "cancel.request").touch()
    try:
        record = json.loads((work / "worker.json").read_text())
        pid = record["pid"]
        if isinstance(pid, int) and pid > 1 and start_time(pid) == record["start_time"]:
            os.kill(pid, signal.SIGTERM)
            deadline = time.monotonic() + 5
            while start_time(pid) == record["start_time"] and time.monotonic() < deadline:
                time.sleep(0.1)
            if start_time(pid) == record["start_time"]:
                return 1
    except (OSError, ValueError, KeyError):
        pass
    return 0


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--workdir", required=True, type=Path)
    parser.add_argument("--assignment", choices=["hw2","hw3","hw4","hw5","hw6","hw7","hw8"], default="hw2")
    parser.add_argument("--cc", default="gcc")
    parser.add_argument("--timeout", type=float, default=120)
    parser.add_argument("--cancel", action="store_true")
    args = parser.parse_args()
    work = args.workdir.resolve()
    if args.cancel:
        return cancel(work)
    artifacts = work / "artifacts"
    artifacts.mkdir(parents=True, exist_ok=True)
    report = dict(assignment=args.assignment, status="environment_error", exit_code=2,
                  score=None, total=None, results=[], environment={}, error=None,
                  cleanup="pending")
    process = None
    deadline = time.monotonic() + args.timeout

    def interrupted(_signum, _frame):
        raise KeyboardInterrupt

    signal.signal(signal.SIGTERM, interrupted)
    signal.signal(signal.SIGINT, interrupted)
    (work / "worker.json").write_text(json.dumps(
        dict(pid=os.getpid(), start_time=start_time(os.getpid()))))

    def run(command, filename, append=False):
        nonlocal process
        if (work / "cancel.request").exists():
            raise KeyboardInterrupt
        path = artifacts / filename
        mode = "ab" if append else "wb"
        with path.open(mode) as log, path.open("rb") as reader:
            if append:
                reader.seek(0, 2)
            process = subprocess.Popen(command, cwd=work / "q2", stdout=log,
                                       stderr=subprocess.STDOUT, start_new_session=True)
            try:
                while True:
                    chunk = reader.read(65536)
                    if chunk:
                        sys.stdout.buffer.write(chunk)
                        sys.stdout.buffer.flush()
                    code = process.poll()
                    if code is not None:
                        rest = reader.read()
                        sys.stdout.buffer.write(rest)
                        sys.stdout.buffer.flush()
                        return code
                    if time.monotonic() >= deadline:
                        raise TimeoutError("Overall build/test deadline exceeded")
                    if (work / "cancel.request").exists():
                        raise KeyboardInterrupt
                    time.sleep(0.05)
            finally:
                terminate_group(process)
                process = None

    try:
        report["environment"] = environment(args.cc)
        print("==> Linux:", report["environment"]["architecture"],
              report["environment"]["kernel"], flush=True)
        print("==> Clean build", flush=True)
        report.update(status="build_failed", exit_code=3)
        code = run(["make", "clean"], "build.log")
        if code == 0:
            code = run(["make", f"CC={args.cc}"], "build.log", append=True)
        report["build_exit_code"] = code
        if code != 0:
            report["error"] = f"Build failed with exit status {code}; see build.log"
        else:
            binary = work / "q2/langar"
            result = artifacts / "judge.json"
            print(f"==> {args.assignment.upper()} public judge", flush=True)
            report.update(status="environment_error", exit_code=2)
            if args.assignment == 'hw2':
                command = [sys.executable, '-B', '-u', str(work/'course/test_langar.py')]
            else:
                command = [sys.executable, '-B', '-u', str(work/'course/test_semester.py'),
                           '--assignment', args.assignment[2:]]
                if os.geteuid() != 0:
                    command = ['sudo', '-n', *command]
            code = run([*command, "--binary", str(binary), "--cc", args.cc,
                        "--json", str(result)], "test.log")
            report["judge_exit_code"] = code
            if code in (0, 1) and result.is_file():
                data = json.loads(result.read_text())
                score, total = data["score"], data["total"]
                cases = data["results"]
                if (not isinstance(score, int) or total != 100 or not 0 <= score <= total
                        or not isinstance(cases, list) or len(cases) != {"hw2":15,"hw3":7,"hw4":6,"hw5":5,"hw6":7,"hw7":6,"hw8":7}[args.assignment]):
                    raise ValueError("Invalid judge result")
                if (code == 0) != all(case["passed"] for case in cases):
                    raise ValueError("Judge exit status and cases disagree")
                report.update(score=score, total=total, results=cases,
                              status="passed" if code == 0 else "test_failed",
                              exit_code=code)
            else:
                report["error"] = f"Judge setup failed (exit {code}); see test.log"
    except KeyboardInterrupt:
        report.update(status="interrupted", exit_code=130, score=None, total=None,
                      error="Run cancelled")
    except TimeoutError as exc:
        status = report["status"]
        report.update(status=status, exit_code=3 if status == "build_failed" else 2,
                      error=str(exc))
    except (OSError, ValueError, KeyError, subprocess.SubprocessError, RuntimeError) as exc:
        report.update(status="environment_error", exit_code=2, score=None, total=None,
                      error=str(exc))
    finally:
        # Do not allow a second termination signal to abort cleanup/report writing.
        signal.signal(signal.SIGTERM, signal.SIG_IGN)
        signal.signal(signal.SIGINT, signal.SIG_IGN)
        terminate_group(process)
        report["cleanup"] = "complete"
        report_path = artifacts / "report.json"
        temporary = report_path.with_suffix(".tmp")
        temporary.write_text(json.dumps(report, ensure_ascii=False, indent=2) + "\n")
        temporary.replace(report_path)
        (work / "worker.json").unlink(missing_ok=True)
    return report["exit_code"]


if __name__ == "__main__":
    sys.exit(main())
