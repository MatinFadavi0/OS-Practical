# Langar Lab — راهنمای فارسی / English guide

[فارسی](#فارسی) · [English](#english)

## فارسی

Langar Lab کد C شما را در لینوکس می‌سازد و آزمون‌های سؤال ۲ را اجرا می‌کند.
در **Linux** به‌صورت پیش‌فرض از سیستم خودتان استفاده می‌کند و VM دانلود
نمی‌کند. در **macOS و Windows** یک ماشین مجازی Ubuntu را با Multipass آماده
و برای اجراهای بعدی دوباره استفاده می‌کند.

صورت سؤال و شرایط تحویل در هندآوت درس است. این فایل فقط راهنمای ابزار است؛
برای اجرای مستقیم داور تمرین ۲ نیز [JUDGE.md](JUDGE.md) را ببینید.

### ۱. پیش‌نیازها

- **Linux، حالت native:** GCC، Make و Python 3.10 یا جدیدتر لازم‌اند.
  روی Ubuntu/Debian فرمان‌های زیر را اجرا کنید. ابزار `sudo` برای آزمون‌های
  دارای namespace و mount از تمرین ۳ به بعد لازم است.

  ```sh
  sudo apt update
  sudo apt install build-essential python3 busybox-static util-linux procps sudo curl ca-certificates git
  ```

- **macOS و Windows:** ابتدا [Multipass را از راهنمای رسمی نصب کنید](https://canonical.com/multipass/docs/latest/how-to-guides/install-multipass/).
  پیش‌نیازهای سیستم و مجازی‌سازی را مطابق همان راهنما بررسی کنید. پس از نصب،
  ترمینال تازه‌ای باز کنید و `multipass version` را اجرا کنید.
  کامپایلر و Python داخل VM نصب می‌شوند؛ روی میزبان Go، Python، Bash یا
  Docker لازم ندارید. نصب‌کنندهٔ macOS/Linux از `sh` و `curl` استفاده می‌کند؛
  نسخهٔ Windows با PowerShell اجرا می‌شود.

برای دریافت مخزن از Git یا گزینهٔ **Code → Download ZIP** در GitHub استفاده کنید.
فرمان‌های بعدی را از ریشهٔ پوشهٔ `OS-Practical` اجرا کنید.

### ۲. نصب ابزار

در macOS یا Linux:

```sh
sh tools/install-langar-lab.sh
./.langar-lab/bin/langar-lab --version
```

در Windows، داخل PowerShell:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File .\tools\install-langar-lab.ps1
.\.langar-lab\bin\langar-lab.exe --version
```

گزینهٔ `ExecutionPolicy` فقط برای همین اجرای نصب‌کننده است و تنظیم دائمی
سیستم را تغییر نمی‌دهد. نصب‌کننده معماری پردازنده را تشخیص می‌دهد، فقط
فایل مناسب را از [نسخهٔ رسمی ابزار](https://github.com/MohammadAminKoohi/OS-Practical/releases/tag/langar-lab-v0.2.0)
می‌گیرد و SHA-256 آن را با فایل داخل مخزن مقایسه می‌کند. فایل اجرایی در
`.langar-lab/bin/` قرار می‌گیرد؛ نیازی به sudo یا تغییر PATH نیست.
اگر نصب قطع شد، همان فرمان را دوباره اجرا کنید. اجرای دوباره ابزار را به
نسخهٔ مشخص‌شده در همین مخزن به‌روز می‌کند و VM یا گزارش‌ها را پاک نمی‌کند.

### ۳. اولین اجرا و آزمون

در macOS/Linux:

```sh
./.langar-lab/bin/langar-lab doctor
./.langar-lab/bin/langar-lab setup
./.langar-lab/bin/langar-lab test hw2 --source ./q2 --json ./hw2-result.json
```

در Windows:

```powershell
.\.langar-lab\bin\langar-lab.exe doctor
.\.langar-lab\bin\langar-lab.exe setup
.\.langar-lab\bin\langar-lab.exe test hw2 --source .\q2 --json .\hw2-result.json
```

در حالت VM، اولین `setup` اینترنت می‌خواهد تا Ubuntu 24.04 و بسته‌ها را
بگیرد. VM پیش‌فرض ۲ CPU، دو گیگابایت RAM و دیسک مجازی ۱۲ گیگابایتی دارد.
در حالت native، `setup` فقط ابزارهای نصب‌شده را بررسی می‌کند.
برای اجراهای بعدی پس از تغییر کد، فقط فرمان `test` را تکرار کنید.

`--source` باید پوشهٔ شامل `Makefile`، `src/` و `include/` را نشان دهد.
مسیر دارای فاصله را داخل کوتیشن بگذارید. کد اولیهٔ تمرین ۲ عمداً ناقص است
و **۵ از ۱۰۰** می‌گیرد؛ این به معنی خراب بودن ابزار نیست.

برای تمرین‌های بعدی، مطابق هندآوت و نسخهٔ منتشرشده، `hw2` را با `hw3` تا
`hw8` جایگزین کنید و کد خودتان را ادامه دهید. خود ابزار پاسخ آماده ندارد.

### ۴. انتخاب محیط اجرا

در مثال‌ها مسیر اجرای macOS/Linux آمده است؛ در Windows از مسیر `.exe` بالا استفاده کنید.

```sh
./.langar-lab/bin/langar-lab test hw2 --source ./q2 --backend native
./.langar-lab/bin/langar-lab setup --backend multipass
./.langar-lab/bin/langar-lab test hw2 --source ./q2 --backend multipass
```

گزینهٔ `native` فقط در Linux کار می‌کند. کاربران Linux نیز می‌توانند
Multipass را نصب و صریحاً حالت VM را انتخاب کنند. در تمرین‌های ۳ تا ۸،
ساخت با کاربر عادی انجام می‌شود ولی داور برای namespace و mount از
`sudo -n` استفاده می‌کند. روی Linux ابتدا `sudo -v` را در همان ترمینال
اجرا کنید؛ اگر تنظیم sudo سیستم اجازه نداد، از Multipass یا روش داوری
مستقیم در راهنمای همان تمرین استفاده کنید. خود CLI را با sudo اجرا نکنید.

### ۵. نتیجه و خطاها

برای هر آزمون، نتیجه و امتیاز چاپ می‌شود. در پایان مسیر گزارش‌ها را می‌بینید؛
آنجا `report.json`، لاگ ساخت و لاگ آزمون، و در صورت تولید، خروجی خام داور
قرار می‌گیرند. `--json` نیز یک نسخهٔ گزارش را در مسیر دلخواه ذخیره می‌کند.

- `0`: همهٔ آزمون‌ها موفق‌اند.
- `1`: آزمون‌ها اجرا شدند؛ یک یا چند مورد شکست خورده است.
- `2`: مشکل محیط، وابستگی، ارتباط VM یا فراخوانی ابزار؛ امتیاز ندارد.
- `3`: خطای ساخت یا پایان مهلت ساخت؛ امتیاز ندارد، `build.log` را ببینید.
- `130`: اجرا با Ctrl-C متوقف شده است؛ امتیاز ندارد.

اگر Multipass پیدا نشد، نصب و فعال بودن سرویس و `multipass version` را
بررسی کنید. برای اتصال کند، `setup --timeout 20m` را امتحان کنید؛ برای
آزمون طولانی‌تر `test hw2 --source ./q2 --timeout 4m`.
اگر checksum نادرست بود، فایل دانلودشده اجرا نمی‌شود؛ مخزن و نصب‌کنندهٔ
همان نسخه را دوباره دریافت کنید. گزارش محلی بازخورد است؛ نمرهٔ نهایی با
اجرای دوبارهٔ داور درس روی کد تحویلی تعیین می‌شود.

### ۶. داخل ابزار چه اتفاقی می‌افتد؟

۱. از `Makefile`، `src/` و `include/` یک کپی می‌گیرد؛ فایل‌های اجرایی و
object میزبان استفاده نمی‌شوند و فایل‌های اصلی شما تغییر نمی‌کنند.
پیوند نمادین در کد منبع پذیرفته نیست؛ خود فایل‌ها را در ساختار خواسته‌شده بگذارید.

۲. روی Linux یا داخل VM، کد کپی‌شده را با Make و GCC می‌سازد و داور همراه
ابزار را اجرا می‌کند. کد برای داوری به سرور دانشگاه یا GitHub فرستاده نمی‌شود؛
در حالت Multipass فقط به ماشین مجازی محلی خودتان منتقل می‌شود.

۳. امتیاز هر آزمون، خطاها، هش کد منبع و مشخصات واقعی محیط Linux را ثبت
می‌کند. گزارش‌ها در پوشهٔ cache کاربر و زیر `langar-lab` می‌مانند و مسیر
دقیق در خروجی چاپ می‌شود. پوشهٔ موقت ساخت پس از اجرا جمع‌آوری می‌شود.

### ۷. مدیریت VM و سؤال ۳

```sh
./.langar-lab/bin/langar-lab status --backend multipass
./.langar-lab/bin/langar-lab shell --backend multipass
./.langar-lab/bin/langar-lab stop --backend multipass
```

`shell` ترمینال داخل VM را باز می‌کند؛ با `exit` برگردید. ابزار فقط VM
اختصاصی خودش را مدیریت می‌کند. `reset --backend multipass` آن VM را حذف و
دوباره ایجاد می‌کند؛ فایل‌های دستی ذخیره‌شده داخل VM از بین می‌روند، ولی
کد و گزارش‌های میزبان باقی می‌مانند. برای اشکال‌زدایی می‌توانید به فرمان
`test` گزینهٔ `--keep-work` بدهید؛ مسیر کپی داخل Linux چاپ می‌شود.

فرمان `test hwN` فقط سؤال ۲ را می‌سنجد. سؤال ۳ شامل مشاهده و توضیح است و
دستی ارزیابی می‌شود. در macOS/Windows با `shell` وارد Linux شوید و مخزن
و نسخهٔ تمرین را **داخل VM** بگیرید؛ پوشهٔ میزبان خودکار mount نمی‌شود.
دستورها و بسته‌های لازم سؤال ۳ در هندآوت و راهنمای همان نسخه آمده‌اند.

## English

Langar Lab builds your C code in Linux and runs the **Q2 public tests**.
**Linux** defaults to native execution, with no VM download. **macOS and
Windows** use a local Ubuntu VM managed by Multipass and reused between runs.
Assignment requirements and submission rules remain in the course handout.
For the direct HW2 judge, see [JUDGE.md](JUDGE.md).

### 1. Install prerequisites

- **Linux native:** GCC, Make and Python 3.10+. On Ubuntu/Debian:

  ```sh
  sudo apt update
  sudo apt install build-essential python3 busybox-static util-linux procps sudo curl ca-certificates git
  ```

  Later namespace/mount tests also need sudo privileges.
- **macOS/Windows:** follow the [official Multipass installation guide](https://canonical.com/multipass/docs/latest/how-to-guides/install-multipass/),
  including its OS and virtualization requirements. Open a new terminal and
  check `multipass version`. The VM receives the compiler and Python; you do
  not need host Go, Python, Bash or Docker to run the executable. The Unix
  installer uses `sh` and `curl`; the Windows installer uses PowerShell.

Clone this repository with Git, or use **Code → Download ZIP** and extract it.
Run the following commands from the `OS-Practical` directory.

### 2. Install and run

macOS/Linux:

```sh
sh tools/install-langar-lab.sh
./.langar-lab/bin/langar-lab doctor
./.langar-lab/bin/langar-lab setup
./.langar-lab/bin/langar-lab test hw2 --source ./q2 --json ./hw2-result.json
```

Windows PowerShell:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File .\tools\install-langar-lab.ps1
.\.langar-lab\bin\langar-lab.exe doctor
.\.langar-lab\bin\langar-lab.exe setup
.\.langar-lab\bin\langar-lab.exe test hw2 --source .\q2 --json .\hw2-result.json
```

The execution-policy option applies only to that installer process. Installers
detect amd64/arm64, download only your executable from
[Langar Lab v0.2.0](https://github.com/MohammadAminKoohi/OS-Practical/releases/tag/langar-lab-v0.2.0),
and verify it against the repository's SHA-256 checksums before running it.
Installation is local to `.langar-lab/bin/`; no administrator access or PATH
change is needed. Rerun the installer to retry or install the version pinned by
your repository checkout; this preserves VMs and reports.

The first VM setup needs internet to download Ubuntu 24.04 and packages.
Defaults are 2 CPUs, 2 GB RAM and a 12 GB virtual disk. Native setup only checks
your installed tools. After editing code, repeat just the `test` command.
`--source` must name the directory containing `Makefile`, `src/` and `include/`;
quote paths containing spaces. The unfinished HW2 starter scoring **5/100** is
expected. Use `hw3` through `hw8` as each assignment is released, continuing your
own earlier code. The helper contains no reference solution.

### 3. Choose native Linux or a VM

These examples use the Unix executable path; use the `.exe` path on Windows.

```sh
./.langar-lab/bin/langar-lab test hw2 --source ./q2 --backend native
./.langar-lab/bin/langar-lab setup --backend multipass
./.langar-lab/bin/langar-lab test hw2 --source ./q2 --backend multipass
```

Native works only on Linux. Linux users can also install Multipass and select
it explicitly. HW3–HW8 compile as your regular user and run privileged judges
through `sudo -n`. On native Linux, run `sudo -v` in the same terminal first;
if your sudo policy cannot authorize the worker, choose Multipass or the direct
judge documented for that release. Run the CLI as your regular user.

### 4. Read results and resolve common problems

Each case prints its result and points. The printed results directory contains
`report.json`, build/test logs and raw judge JSON when produced. `--json` saves
an extra report at your chosen path. Exit codes are:

- `0`: all tests passed.
- `1`: tests completed with failing cases.
- `2`: usage, dependency, environment or VM transport error; no score.
- `3`: build failed or timed out; no score. Read `build.log`.
- `130`: interrupted with Ctrl-C; no score.

If Multipass is unavailable, check installation, its running service and
`multipass version`. For slow setup, use `setup --timeout 20m`; for a longer test,
use `test hw2 --source ./q2 --timeout 4m`. A checksum mismatch prevents
installation: obtain the checkout and installer for the same release again.
Local scores are feedback; final grading reruns the course-owned judge on your
submission.

### 5. How it works

1. Copies `Makefile`, `src/` and `include/` into a source snapshot. Your original
   files stay unchanged; host executables and object files are not used. Source
   symlinks are rejected, so copy actual files into the required layout.
2. Builds the copy with Make/GCC inside Linux and runs its embedded public
   judge. Code is processed locally: Multipass transfers it to your own local
   VM, with no upload to GitHub or a university grading service.
3. Saves case results, source hash and actual Linux environment details. Host
   reports remain under the user-cache directory in `langar-lab`, at the path
   printed after each run. Temporary Linux build workspaces are cleaned up.

### 6. VM commands and Q3

```sh
./.langar-lab/bin/langar-lab status --backend multipass
./.langar-lab/bin/langar-lab shell --backend multipass
./.langar-lab/bin/langar-lab stop --backend multipass
```

Use `exit` to leave the guest shell. The helper manages only its own course VM.
`reset --backend multipass` deletes and recreates that VM, including files you
saved manually inside it; host code and reports remain. Add `--keep-work` to a
test to retain its Linux copy for debugging; the workspace path is printed.

`test hwN` evaluates Q2 only. Q3 observations/explanations are manually graded.
On macOS/Windows, enter the VM with `shell` and obtain the repository/release
**inside Linux**; the host repository is not automatically mounted. Follow that
release's handout and judge guide for Q3 commands and packages.

### Optional: manual download or build from source / دریافت دستی یا ساخت از منبع

Executables and `SHA256SUMS` are on the
[release page](https://github.com/MohammadAminKoohi/OS-Practical/releases/tag/langar-lab-v0.2.0).
Use `darwin` for macOS, `linux` for Linux and `windows` for Windows; use `amd64`
for Intel/AMD 64-bit and `arm64` for Apple Silicon/ARM64. Verify SHA-256 before
running. On Unix, run `chmod +x FILE`. Binaries are unsigned; if your OS blocks
execution, building from source is an alternative.

برای دریافت دستی، فایل مناسب سیستم و معماری خود را از صفحهٔ نسخه بگیرید
و SHA-256 آن را با فایل داخل مخزن مقایسه کنید. فایل‌های اجرایی امضای سیستم‌عامل
ندارند؛ در صورت محدودیت اجرا می‌توانید ابزار را از منبع بسازید.
ساخت از منبع اختیاری است و به [Go 1.22 یا جدیدتر](https://go.dev/doc/install)
نیاز دارد؛ وابستگی Go دیگری وجود ندارد.

With [Go 1.22+](https://go.dev/doc/install), no third-party Go modules are needed.
Run inside `tools/langar-lab/` (the Windows output uses `.exe`):

```sh
go test ./...
go build -trimpath -o dist/langar-lab ./cmd/langar-lab
```

On Windows use `-o dist/langar-lab.exe`. With Python 3 installed,
`python3 scripts/build.py` cross-builds all six targets and generates checksums.
Build output goes in `dist/`; it is not part of your assignment submission.
Cross-building does not imply runtime testing on every host.
