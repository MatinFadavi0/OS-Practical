# Operating Systems Practical Resources

Student starter code, public tests, and judge tools for the Operating Systems course.
Questions and submission requirements are provided separately in the course handout.

این مخزن شامل کد اولیه، آزمون‌ها و ابزارهای داوری است.
صورت سؤال‌ها و الزامات تحویل جداگانه در هندآوت درس ارائه می‌شوند.

## HW02 — `hw02-v1.0.1`

- [Langar starter](assignments/hw02/q2/)
- [آماده‌سازی لینوکس](docs/setup.fa.md) / [Linux setup](docs/setup.en.md)
- [راهنمای آزمون‌ها](tests/hw02/README.fa.md) / [Test runner guide](tests/hw02/README.en.md)

```sh
make -C assignments/hw02/q2 CC=gcc
python3 tests/hw02/test_langar.py
```

The starter builds but intentionally leaves `run_command` unimplemented. Only
the CLI test passes initially (5/100 functional points).

برای بررسی ساختار پوشهٔ تحویل / To check an extracted submission directory:

```sh
python3 tools/check_submission.py /path/to/extracted-submission
```

Follow the handout for the required files and videos. Keep your own implementation
as the project grows; use the release tag announced for each assignment.
