# 🖼️ imageopt — بهینه‌ساز حرفه‌ای تصاویر برای وب

<div dir="rtl">

[![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![Dependency](https://img.shields.io/badge/libvips-required-orange)](https://libvips.github.io/libvips/)

ماژول `imageopt` یک کتابخانه‌ی Go است که تصاویر رستر (JPEG، PNG، TIFF و ...) را به‌صورت هوشمند به فرمت **WebP** تبدیل و بهینه می‌کند. این ماژول بر پایه‌ی کتابخانه‌ی قدرتمند [`bimg`](https://github.com/h2non/bimg) (که خود از `libvips` استفاده می‌کند) ساخته شده و برای استقرار در محیط‌های Production طراحی شده است.

---

## 📋 فهرست مطالب

- [ویژگی‌ها](#-ویژگی‌ها)
- [معماری و منطق](#-معماری-و-منطق)
- [پیش‌نیازها](#-پیش‌نیازها)
- [نصب](#-نصب)
- [استفاده‌ی سریع](#-استفاده‌ی-سریع)
- [مثال‌های کامل](#-مثال‌های-کامل)
- [جزئیات الگوریتم](#-جزئیات-الگوریتم)
- [خروجی لاگ](#-خروجی-لاگ)
- [نکات Production](#-نکات-production)
- [محدودیت‌ها](#-محدودیت‌ها)

---

## ✨ ویژگی‌ها

- **تبدیل خودکار به WebP** — خروجی همیشه WebP است که نسبت به JPEG تا ۳۰٪ و نسبت به PNG تا ۷۰٪ کوچک‌تر است
- **کاهش رزولوشن هوشمند** — بر اساس عرض اصلی تصویر، در صورت نیاز عرض را به مقادیر استاندارد وب کاهش می‌دهد
- **کیفیت تطبیقی** — کیفیت فشرده‌سازی را بر اساس چگالی اطلاعات (bytes per megapixel) محاسبه می‌کند
- **شارپ‌کردن پس از resize** — اگر تصویر کوچک شود، به‌صورت خودکار sharpen اعمال می‌شود
- **حذف متادیتا** — تمام اطلاعات Exif، GPS و ICC به‌صورت خودکار پاک می‌شوند
- **Interlace / Progressive** — فعال‌سازی پیش‌فرض برای بارگذاری تدریجی در مرورگر
- **بدون تخصیص حافظه‌ی اضافی** — پردازش مستقیم روی بافر خام

---

## 🏗️ معماری و منطق

```
inputFile ──► ReadFile ──► bimg.NewImage
                                │
                          Size() → width, height, pixels, fileSize
                                │
                    ┌───────────┴────────────┐
                    │                        │
           calculateTargetWidth        calculateQuality
           (بر اساس عرض اصلی)       (بر اساس bytes/MP)
                    │                        │
                    └───────────┬────────────┘
                                │
                       bimg.Options{...}
                       + Sharpen (در صورت resize)
                                │
                          img.Process()
                                │
                         WriteFile ──► outputFile
                                │
                          log.Printf(stats)
```

---

## 📦 پیش‌نیازها

### ۱. نصب `libvips`

`bimg` به کتابخانه‌ی سیستمی `libvips` نیاز دارد. نسخه‌ی **8.8+** توصیه می‌شود.

**Ubuntu / Debian:**
```bash
apt-get update && apt-get install -y libvips-dev
```

**macOS (Homebrew):**
```bash
brew install vips
```

**Alpine Linux (Docker):**
```bash
apk add vips-dev build-base
```

**CentOS / RHEL:**
```bash
yum install -y vips-devel
```

### ۲. Go

نسخه‌ی **1.21** یا بالاتر لازم است.

```bash
go version  # باید go1.21+ نشان دهد
```

---

## 🚀 نصب

```bash
go get github.com/your-username/imageopt
```

یا اگر ماژول را به‌صورت لوکال دارید، در `go.mod` پروژه‌تان:

```
require github.com/your-username/imageopt v1.0.0
```

سپس وابستگی‌ها را resolve کنید:

```bash
go mod tidy
```

---

## ⚡ استفاده‌ی سریع

```go
package main

import (
    "log"
    "github.com/your-username/imageopt"
)

func main() {
    err := imageopt.OptimizeImageForWeb("photo.jpg", "photo.webp")
    if err != nil {
        log.Fatal(err)
    }
    // 2024/01/15 10:23:45 Done | 4000x3000 → 2000px | 8192KB → 1241KB | Q=48
}
```

تنها یک تابع عمومی وجود دارد:

```go
func OptimizeImageForWeb(inputFile string, outputFile string) error
```

| پارامتر | نوع | توضیح |
|---------|-----|-------|
| `inputFile` | `string` | مسیر کامل فایل ورودی (JPEG، PNG، TIFF، ...) |
| `outputFile` | `string` | مسیر فایل خروجی `.webp` |
| **بازگشتی** | `error` | در صورت موفقیت `nil`، در غیر این صورت خطا |

---

## 📚 مثال‌های کامل

### مثال ۱ — بهینه‌سازی یک تصویر ساده

```go
package main

import (
    "fmt"
    "log"
    "github.com/your-username/imageopt"
)

func main() {
    input  := "original.png"
    output := "optimized.webp"

    if err := imageopt.OptimizeImageForWeb(input, output); err != nil {
        log.Fatalf("خطا در بهینه‌سازی: %v", err)
    }

    fmt.Println("✅ تصویر با موفقیت بهینه شد.")
}
```

**خروجی لاگ:**
```
2024/01/15 10:23:45 Done | 1920x1080 → 1920px | 512KB → 87KB | Q=76
```

---

### مثال ۲ — پردازش دسته‌ای (Batch Processing) یک پوشه

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"

    "github.com/your-username/imageopt"
)

func main() {
    inputDir  := "./images/original"
    outputDir := "./images/optimized"

    // ساخت پوشه‌ی خروجی اگر وجود نداشته باشد
    if err := os.MkdirAll(outputDir, 0755); err != nil {
        log.Fatal(err)
    }

    entries, err := os.ReadDir(inputDir)
    if err != nil {
        log.Fatal(err)
    }

    var success, failed int

    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }

        name := entry.Name()
        ext  := strings.ToLower(filepath.Ext(name))

        // فقط فرمت‌های پشتیبانی‌شده
        switch ext {
        case ".jpg", ".jpeg", ".png", ".tiff", ".webp":
            // ادامه بده
        default:
            fmt.Printf("رد شد: %s (فرمت پشتیبانی‌نشده)\n", name)
            continue
        }

        inputPath  := filepath.Join(inputDir, name)
        outputName := strings.TrimSuffix(name, ext) + ".webp"
        outputPath := filepath.Join(outputDir, outputName)

        if err := imageopt.OptimizeImageForWeb(inputPath, outputPath); err != nil {
            log.Printf("❌ خطا در پردازش %s: %v", name, err)
            failed++
            continue
        }

        success++
    }

    fmt.Printf("\n✅ موفق: %d  |  ❌ ناموفق: %d\n", success, failed)
}
```

---

### مثال ۳ — پردازش موازی با Worker Pool

برای پردازش سریع‌تر چندین تصویر به‌صورت همزمان:

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"
    "sync"

    "github.com/your-username/imageopt"
)

type job struct {
    input  string
    output string
}

func main() {
    inputDir  := "./raw"
    outputDir := "./web"
    workers   := 4 // تعداد goroutine موازی

    os.MkdirAll(outputDir, 0755)

    entries, _ := os.ReadDir(inputDir)

    jobs := make(chan job, len(entries))
    var wg sync.WaitGroup

    // راه‌اندازی worker pool
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for j := range jobs {
                if err := imageopt.OptimizeImageForWeb(j.input, j.output); err != nil {
                    log.Printf("خطا: %s → %v", j.input, err)
                }
            }
        }()
    }

    // ارسال کارها به کانال
    for _, e := range entries {
        ext := strings.ToLower(filepath.Ext(e.Name()))
        if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
            in  := filepath.Join(inputDir, e.Name())
            out := filepath.Join(outputDir, strings.TrimSuffix(e.Name(), ext)+".webp")
            jobs <- job{input: in, output: out}
        }
    }

    close(jobs)
    wg.Wait()

    fmt.Println("✅ پردازش موازی به پایان رسید.")
}
```

> **توجه:** تعداد workers را با توجه به تعداد هسته‌های CPU و RAM موجود تنظیم کنید. پردازش تصاویر ۸K ممکن است روی هر هسته ۵۰۰MB RAM مصرف کند.

---

### مثال ۴ — استفاده در HTTP Handler (آپلود و بهینه‌سازی آنی)

```go
package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"

    "github.com/your-username/imageopt"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "فقط POST مجاز است", http.StatusMethodNotAllowed)
        return
    }

    // محدودیت حجم آپلود: 50MB
    r.Body = http.MaxBytesReader(w, r.Body, 50<<20)

    // دریافت فایل از فرم
    file, header, err := r.FormFile("image")
    if err != nil {
        http.Error(w, "خطا در دریافت فایل", http.StatusBadRequest)
        return
    }
    defer file.Close()

    // ذخیره‌ی موقت فایل ورودی
    tmpIn := filepath.Join(os.TempDir(), "upload_"+header.Filename)
    out, _ := os.Create(tmpIn)
    io.Copy(out, file)
    out.Close()
    defer os.Remove(tmpIn) // پاکسازی پس از اتمام

    // مسیر خروجی WebP
    tmpOut := tmpIn + ".webp"
    defer os.Remove(tmpOut)

    // بهینه‌سازی
    if err := imageopt.OptimizeImageForWeb(tmpIn, tmpOut); err != nil {
        http.Error(w, fmt.Sprintf("خطا در بهینه‌سازی: %v", err), http.StatusInternalServerError)
        return
    }

    // ارسال فایل WebP به کلاینت
    w.Header().Set("Content-Type", "image/webp")
    w.Header().Set("Content-Disposition", `attachment; filename="optimized.webp"`)
    http.ServeFile(w, r, tmpOut)
}

func main() {
    http.HandleFunc("/upload", uploadHandler)
    fmt.Println("سرور در حال اجرا روی :8080")
    http.ListenAndServe(":8080", nil)
}
```

**تست با curl:**
```bash
curl -X POST http://localhost:8080/upload \
  -F "image=@/path/to/photo.jpg" \
  -o result.webp
```

---

### مثال ۵ — نمونه Dockerfile برای محیط Production

```dockerfile
# ─── Build Stage ─────────────────────────────────────────
FROM golang:1.22-bookworm AS builder

RUN apt-get update && \
    apt-get install -y libvips-dev pkg-config && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o optimizer ./cmd/optimizer

# ─── Runtime Stage ────────────────────────────────────────
FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y libvips42 && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/optimizer /usr/local/bin/optimizer

ENTRYPOINT ["optimizer"]
```

```bash
# Build
docker build -t imageopt:latest .

# اجرا با mount پوشه‌ی تصاویر
docker run --rm \
  -v /host/images:/images \
  imageopt:latest /images/input.jpg /images/output.webp
```

> **مهم:** چون `bimg` از CGO استفاده می‌کند، حتماً `CGO_ENABLED=1` را هنگام build تنظیم کنید.

---

## 🔬 جزئیات الگوریتم

### محاسبه‌ی عرض هدف (`calculateTargetWidth`)

این تابع تصاویر با رزولوشن بسیار بالا را به عرض‌های استاندارد وب کاهش می‌دهد:

| عرض اصلی تصویر | عرض خروجی | توضیح |
|----------------|-----------|-------|
| ≤ 1200px | بدون تغییر | تصویر کوچک، نیازی به کاهش نیست |
| 1200px – 2199px | بدون تغییر | Large Web، مناسب است |
| 2200px – 2799px | **1600px** | تصویر 4K |
| 2800px – 3599px | **2000px** | تصویر 5K |
| ≥ 3600px | **2000px** | تصویر 6K–8K |

> تصاویری که عرضشان کمتر از ۱۲۰۱ پیکسل است **بدون هیچ تغییری** در رزولوشن پردازش می‌شوند.

---

### محاسبه‌ی کیفیت (`calculateQuality`)

کیفیت پایه `76` است. بر اساس دو معیار به‌صورت تجمعی کاهش می‌یابد:

**۱. چگالی اطلاعات (bytes per megapixel) — تصاویر پر نویز:**

| چگالی | کاهش کیفیت | مثال |
|--------|------------|------|
| > 8 MB/MP | −28 | عکس RAW بسیار پرجزئیات |
| > 5 MB/MP | −22 | عکس با کیفیت بالا |
| > 3.5 MB/MP | −15 | تصویر متوسط |
| > 2.5 MB/MP | −8 | تصویر نسبتاً سبک |
| ≤ 2.5 MB/MP | 0 | تصویر بهینه یا کوچک |

**۲. رزولوشن بسیار بالا — جریمه‌ی اضافی:**

| مگاپیکسل | کاهش اضافی |
|----------|-----------|
| > 24 MP (~8K) | −10 |
| > 16 MP (~6K) | −6 |
| > 12 MP (~5K) | −4 |
| ≤ 12 MP | 0 |

**محدوده‌ی نهایی:** کیفیت همیشه بین `45` و `76` باقی می‌ماند (clamp).

**مثال محاسبه:**
```
تصویر: 7680×4320 (8K) | حجم: 40MB
megapixels = 33.2 MP
bytesPerMP  = 40,000,000 / 33.2 = ~1,204,819  →  کاهش 0
megapixels > 24                               →  کاهش −10
کیفیت نهایی = 76 − 10 = 66
```

---

## 📝 خروجی لاگ

هر بار که `OptimizeImageForWeb` با موفقیت اجرا شود، یک خط لاگ چاپ می‌کند:

```
2024/01/15 10:23:45 Done | 4000x3000 → 2000px | 8192KB → 1241KB | Q=48
```

| بخش | توضیح |
|-----|-------|
| `4000x3000` | ابعاد اصلی تصویر (عرض × ارتفاع) |
| `→ 2000px` | عرض خروجی پس از resize |
| `8192KB` | حجم فایل ورودی |
| `→ 1241KB` | حجم فایل خروجی |
| `Q=48` | کیفیت WebP انتخاب‌شده |

برای capture کردن لاگ‌ها در production:
```bash
./optimizer 2>&1 | tee /var/log/imageopt.log
```

---

## 🏭 نکات Production

**۱. مدیریت خطا:** همیشه خطای برگشتی را بررسی کنید. اگر `libvips` نصب نباشد یا فایل ورودی خراب باشد، error برمی‌گردد.

**۲. فایل‌های موقت:** در سرور، از `/tmp` برای فایل‌های ورودی موقت استفاده کنید و پس از پردازش آن‌ها را حذف کنید (مثال ۴ را ببینید).

**۳. مدیریت حافظه:** پردازش تصاویر ۸K ممکن است روی هر درخواست ۵۰۰MB+ RAM مصرف کند. اگر load بالاست، تعداد worker را محدود کنید یا از memory limit در Docker استفاده کنید:
```bash
docker run --memory=2g imageopt:latest
```

**۴. CGO و Cross-compile:** چون `bimg` از CGO استفاده می‌کند، cross-compile مستقیم ممکن نیست. برای ساخت باینری لینوکس روی مک از Docker استفاده کنید.

**۵. مانیتورینگ:** می‌توانید لاگ‌ها را parse کنید تا نرخ compression را به Prometheus یا Datadog بفرستید:
```bash
# میانگین نرخ compression
grep "Done" imageopt.log | awk -F'[| ]' '{gsub("KB",""); print $6/$8}' | awk '{s+=$1;c++} END {print s/c}'
```

**۶. ایمنی در concurrent استفاده:** `OptimizeImageForWeb` stateless است و می‌توان آن را به‌صورت موازی از چندین goroutine فراخواند.

---

## ⚠️ محدودیت‌ها

- فرمت خروجی **همیشه WebP** است؛ در حال حاضر AVIF یا JPEG پشتیبانی نمی‌شود
- فایل‌های ورودی باید از فایل‌سیستم خوانده شوند؛ پشتیبانی از `io.Reader` در نسخه‌های بعدی
- در محیط‌هایی که `libvips` در دسترس نیست (مثل AWS Lambda پیش‌فرض) نیاز به layer سفارشی دارید
- تصاویر GIF انیمیشن‌دار به درستی پردازش نمی‌شوند (تنها فریم اول حفظ می‌شود)
- بهینه‌سازی ارتفاع تصویر انجام نمی‌شود؛ فقط عرض تنظیم می‌گردد

---

## 📄 لایسنس

MIT License — آزادانه استفاده، تغییر، و توزیع کنید.

</div>
