# Operating Systems Practical Resources

Student starter code, public tests, and tools for the Operating Systems course.
The course LaTeX handout is the source for the questions, requirements, and
submission rules. This repository does not duplicate the problem statements in Markdown.

منبع صورت سؤال‌ها، الزامات و روش تحویل، هندآوت درس با منبع LaTeX است.
این مخزن شامل کد اولیه، آزمون‌ها و ابزارهای مورد نیاز دانشجویان است.

## HW02 — `hw02-v1.0.0`

- [Langar starter](assignments/hw02/q2/)
- [Q3 answer template](assignments/hw02/q3/answers.md)
- [LaTeX handout fragments](assignments/hw02/handouts/)
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
