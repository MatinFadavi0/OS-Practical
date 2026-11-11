<div dir="rtl" align="right">

# پاسخ تمرین ۳ — لنگر: جداسازی پردازه و نام میزبان

## محیط اجرا

- کرنل: <bdi>`Linux 7.2.11-arch1-1`</bdi>
- معماری: <bdi>`x86_64`</bdi>

## پرسش‌های گزارش

### ۱. چرا پس از <bdi>`unshare`</bdi> مربوط به <bdi>PID</bdi> هنوز به <bdi>`fork`</bdi> نیاز داریم؟

<bdi>PID</bdi> یک پردازه و فضای نام <bdi>PID</bdi> آن در تمام عمرش ثابت است. برای همین <bdi>`unshare(CLONE_NEWPID)`</bdi> **خودِ پردازهٔ صدازننده را جابه‌جا نمی‌کند**، فقط فرزندان **بعدی** آن را در فضای جدید می‌سازد. اولین فرزندی که بعد از آن با <bdi>`fork`</bdi> ساخته شود، در فضای جدید <bdi>PID 1</bdi> می‌گیرد. اگر بعد از <bdi>`unshare`</bdi> مستقیم <bdi>`execvp`</bdi> می‌کردیم، برنامه با همان <bdi>PID</bdi> قبلی اجرا می‌شد.

ترتیب واقعی فراخوان‌ها در <bdi>`src/run.c`</bdi>:

| خط | پردازه | فراخوان |
|---|---|---|
| ۴۲ | پدر (لنگر) | <bdi>`unshare(CLONE_NEWPID)`</bdi> |
| ۴۸ | پدر | <bdi>`fork()`</bdi> |
| ۲۰ | فرزند | <bdi>`unshare(CLONE_NEWUTS)`</bdi> |
| ۲۷ | فرزند | <bdi>`sethostname(...)`</bdi> |
| ۶۴ | فرزند | <bdi>`execvp(argv[0], argv)`</bdi> |
| ۷۹ | پدر | <bdi>`waitpid(pid, ...)`</bdi> |

پدر بیرون از فضای جدید می‌ماند، پس فرزند <bdi>`PPID = 0`</bdi> می‌بیند:

<div dir="ltr" align="left">

```
$ sudo ./langar run --isolate -- /bin/sh -c 'echo pid=$$ ppid=$PPID'
pid=1 ppid=0
```

</div>

### ۲. کدام اطلاعات این هفته جدا شده‌اند و کدام منابع هنوز مشترک‌اند؟

**جدا شده:**

- **شمارهٔ پردازه‌ها (<bdi>PID namespace</bdi>):** برنامه خودش را <bdi>PID 1</bdi> می‌بیند.
- **اسم میزبان (<bdi>UTS namespace</bdi>):** اسم داخل لنگر <bdi>`langar`</bdi> یا مقدار <bdi>`--hostname`</bdi> است و اسم سیستم اصلی عوض نمی‌شود.

**هنوز مشترک:**

- فایل‌سیستم و <bdi>mount</bdi>ها، از جمله <bdi>`/proc`</bdi>
- شبکه
- کاربرها (<bdi>UID</bdi>ها)
- ارتباط بین پردازه‌ای (<bdi>IPC</bdi>)
- منابع <bdi>CPU</bdi> و حافظه

## محدودیت‌های شناخته‌شده

- فضای نام <bdi>mount</bdi> و <bdi>`/proc`</bdi> جدا ندارد، پس <bdi>`ps`</bdi> داخل لنگر هنوز پردازه‌های سیستم اصلی را نشان می‌دهد.
- سیگنال‌ها به <bdi>PID 1</bdi> ایزوله منتقل یا مدیریت نمی‌شوند.
- بدون دسترسی <bdi>root</bdi>، <bdi>`--isolate`</bdi> با خطای <bdi>`Operation not permitted`</bdi> و کد 1 تمام می‌شود.

## زمان شروع پاسخ پرسش‌های ویدیو

| زمان | پرسش |
|------|-----|
| <bdi>`5:06`</bdi> | ۱. چرا <bdi>`getpid`</bdi> عدد ۱ می‌دهد ولی <bdi>`ps`</bdi> با <bdi>`/proc`</bdi> موروثی نمای دیگری نشان می‌دهد؟ |
| <bdi>`5:47`</bdi> | ۲. اگر <bdi>`sethostname`</bdi> پیش از ایجاد <bdi>UTS namespace</bdi> اجرا شود چه می‌شود و کد من چطور جلویش را می‌گیرد؟ |

</div>
