# راهنمای اجرای داور / Running the judge

این مخزن فقط کد اولیه و ابزارهای داوری را دارد. صورت سؤال و شرایط تحویل را
از هندآوت درس دنبال کنید. پس از دریافت مخزن و نصب ابزارها، ساخت برنامه و
اجرای داور بدون sudo و اتصال شبکه انجام می‌شوند.

## پیش‌نیازها

در لینوکس به GCC، Make و Python 3.10 یا جدیدتر نیاز دارید. ابزارهای داوری
فقط از کتابخانهٔ استاندارد Python استفاده می‌کنند و بستهٔ pip لازم ندارند.
روی Debian/Ubuntu:

```sh
sudo apt update
sudo apt install build-essential python3 git
```

## دریافت و ساخت

```sh
git clone https://github.com/MohammadAminKoohi/OS-Practical.git
cd OS-Practical
make -C q2 clean
make -C q2 CC=gcc
```

خروجی ساخت، فایل `q2/langar` است. کد خود را در `q2/src/run.c` تکمیل کنید.

## اجرای آزمون‌ها

از ریشهٔ مخزن اجرا کنید:

```sh
python3 judge/test_langar.py
```

برای هر آزمون، `PASS` یا `FAIL`، امتیاز آن و در صورت شکست، توضیح خطا چاپ
می‌شود. مجموع امتیاز عملکردی ۱۰۰ است. کد اولیه عمداً ناقص است: فقط آزمون
خط فرمان را می‌گذراند و امتیاز اولیه **۵ از ۱۰۰** طبیعی است.

وزن `direct_child` برابر ۱۵ است؛ `argument_boundaries`، `parent_waits`
و `normal_exit_status` هرکدام ۱۰ امتیاز و یازده آزمون دیگر هرکدام ۵ امتیاز
دارند. امتیاز هر مورد با عبور از تمام بررسی‌های همان مورد داده می‌شود.

کد خروج ابزار:

- `0`: همهٔ آزمون‌ها موفق بودند.
- `1`: یک یا چند آزمون عملکردی ناموفق بود.
- `2`: آماده‌سازی یا فراخوانی ابزار اشکال داشت؛ مانند نبودن فایل اجرایی یا کامپایلر.

داور برنامهٔ کمکی `probe.c` را در پوشهٔ موقت کامپایل می‌کند؛ آن را دستی
نسازید و در کد خود استفاده نکنید. هر اجرای آزمون مهلت پنج‌ثانیه‌ای دارد.

## ذخیرهٔ خروجی و نتیجهٔ JSON

برای ثبت هر دو خروجی و نگه‌داشتن کد خروج واقعی آزمون:

```sh
python3 judge/test_langar.py > q2/test-results.txt 2>&1
status=$?
cat q2/test-results.txt
echo "$status"
```

برای ذخیرهٔ نتیجهٔ ماشینی یا آزمودن یک فایل اجرایی دیگر:

```sh
python3 judge/test_langar.py --json /tmp/langar-result.json
python3 judge/test_langar.py --binary /path/to/submission/q2/langar --cc gcc
```

JSON شامل `score`، `total` و برای هر آزمون، `id`، `points`، `earned`،
`passed` و `error` است. مسیر `--binary` می‌تواند به برنامه‌ای خارج از مخزن
اشاره کند؛ آزمون‌های همین مخزن روی آن اجرا می‌شوند.

## بررسی ساختار تحویل

پس از آماده کردن تمام فایل‌های تحویل مطابق هندآوت، مسیر پوشهٔ استخراج‌شده
را بدهید؛ داخل آن باید پوشه‌های `q2/` و `q3/` قرار داشته باشند:

```sh
python3 judge/check_submission.py /path/to/extracted-submission
```

این ابزار مسیرها، وجود فایل‌های غیرخالی، UTF-8 بودن متن و سرآیند اولیهٔ
ویدیوها را بررسی می‌کند؛ جایگزین آزمون عملکردی یا بررسی محتوای گزارش و
ویدیو نیست. پیش از بررسی ساختار تحویل، فایل‌های ساخته‌شده را پاک کنید:

```sh
make -C /path/to/extracted-submission/q2 clean
```

برای دریافت JSON از این بررسی نیز می‌توانید `--json /tmp/layout.json`
اضافه کنید. خروج با `0` یعنی موفقیت بررسی‌ها، `1` یعنی اشکال در تحویل و
`2` یعنی استفادهٔ نادرست از ابزار.

## English quick reference

Use Linux with GCC, Make, and Python 3.10+. Run commands from the repository root:

```sh
make -C q2 clean
make -C q2 CC=gcc
python3 judge/test_langar.py
```

The untouched starter builds and scores **5/100**; implement `q2/src/run.c`.
Each test prints its result and points. Functional test exit codes are 0 for all
passed, 1 for failed tests, and 2 for setup/usage errors. No pip packages are needed.

Use `--binary /path/to/langar` to test another executable, `--cc gcc` to select
the probe compiler, and `--json /tmp/result.json` to save per-test scores.
Capture both output streams in `q2/test-results.txt` using the command above.

After cleaning the built executable and objects, run
`python3 judge/check_submission.py /path/to/extracted-submission` to check the
required submission files. This checks layout and basic file formats only;
follow the separately provided course handout for requirements.
