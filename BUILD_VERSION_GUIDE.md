# Build Version Integration Guide

این راهنما نحوه استفاده از سیستم build version در ArasAuth را توضیح می‌دهد.

## تغییرات اعمال شده

### 1. Dockerfile
- اضافه شدن ARG declarations برای `BUILD_VERSION`, `BUILD_TIME`, `GIT_COMMIT`
- استفاده از این ARGها در ldflags بجای مقادیر hardcoded

### 2. Makefile
- اضافه شدن target جدید `docker-build-versioned` برای build با اطلاعات نسخه خودکار

### 3. scripts/setup.sh
- اضافه شدن function `generate_build_info()` برای تولید خودکار اطلاعات build
- اضافه شدن command جدید `docker-build` برای build با version info
- به‌روزرسانی `setup_environment()` برای اضافه کردن build info به .env

## روش‌های استفاده

### روش 1: استفاده از Makefile (توصیه شده)
```bash
# Build با اطلاعات نسخه خودکار از git
make docker-build-versioned
```

### روش 2: استفاده از setup script
```bash
# Build با اطلاعات نسخه خودکار
./scripts/setup.sh docker-build
```

### روش 3: استفاده از docker-compose
```bash
# ابتدا فایل .env را تنظیم کنید
echo "BUILD_VERSION=1.1.0" >> .env
echo "BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)" >> .env
echo "GIT_COMMIT=$(git rev-parse --short HEAD)" >> .env

# سپس build کنید
docker-compose build
```

### روش 4: Build مستقیم با docker
```bash
docker build \
  --build-arg BUILD_VERSION=1.1.0 \
  --build-arg BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
  --build-arg GIT_COMMIT=$(git rev-parse --short HEAD) \
  -t aras-auth:latest .
```

## تست و تأیید

برای تست اینکه اطلاعات نسخه درست در ایمیج قرار گرفته:

```bash
docker run --rm aras-auth:latest ./aras-auth --version
```

خروجی باید شامل version، build time و git commit باشد.

## مزایا

1. **اطلاعات نسخه درون ایمیج**: مقادیر در زمان build در ایمیج قرار می‌گیرند
2. **انعطاف‌پذیری**: امکان تنظیم مقادیر از طریق ARGها
3. **خودکارسازی**: امکان تولید خودکار مقادیر از git
4. **سازگاری**: حفظ سازگاری با روش‌های مختلف build
