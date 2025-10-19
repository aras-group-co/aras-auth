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

## Push به GitHub Container Registry

### Authentication

قبل از push، باید با GitHub Container Registry احراز هویت کنید:

```bash
# Login to GitHub Container Registry
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin

# یا با username/password
docker login ghcr.io
```

### Push کردن ایمیج

```bash
# Push ایمیج latest
make docker-push

# Push ایمیج با version info
make docker-push-versioned
```

### استفاده از ایمیج در docker-compose

حالا می‌توانید از ایمیج GHCR در `docker-compose.yml` استفاده کنید:

```yaml
services:
  aras_auth:
    image: ghcr.io/aras-group-co/aras-auth:latest
    # یا برای version مشخص:
    # image: ghcr.io/aras-group-co/aras-auth:v1.1.0
```

## Version Management Best Practices

### چرا از `latest` در Production استفاده نکنیم؟

❌ **مشکلات `latest` در Production:**
- **غیرقابل پیش‌بینی**: نمی‌دانید دقیقاً کدام version را اجرا می‌کنید
- **عدم امکان Rollback**: اگر مشکلی پیش آید، نمی‌توانید به version قبلی برگردید
- **ناسازگاری بین سرورها**: سرورهای مختلف ممکن است version های متفاوت `latest` داشته باشند
- **مشکلات Debugging**: نمی‌توانید دقیقاً بگویید کدام version مشکل دارد
- **Breaking Changes**: اگر version جدیدی push شود، ممکن است سیستم شما خراب شود

### استراتژی Version Management پیاده‌سازی شده

✅ **Production (docker-compose.yml):**
```yaml
services:
  aras_auth:
    image: ghcr.io/aras-group-co/aras-auth:${ARAS_AUTH_VERSION:-v1.1.0}
```
- **Default**: همیشه از version مشخص استفاده می‌کند (v1.1.0)
- **Safe**: اگر متغیر set نشود، از version stable استفاده می‌کند

✅ **Development (docker-compose.dev.yml):**
```yaml
services:
  aras_auth:
    image: ghcr.io/aras-group-co/aras-auth:${ARAS_AUTH_VERSION:-latest}
```
- **Default**: از `latest` استفاده می‌کند (برای تست feature های جدید)
- **انعطاف‌پذیر**: می‌توان version مشخص را تست کرد

### نحوه استفاده

**Production (بدون تنظیم متغیر):**
```bash
docker-compose up
# استفاده می‌کند از: ghcr.io/aras-group-co/aras-auth:v1.1.0
```

**Development (با latest):**
```bash
# استفاده از default (latest)
docker-compose -f docker-compose.dev.yml up

# یا صراحتاً latest
export ARAS_AUTH_VERSION=latest
docker-compose up
```

**Testing version خاص:**
```bash
export ARAS_AUTH_VERSION=v1.0.5
docker-compose up
# استفاده می‌کند از: ghcr.io/aras-group-co/aras-auth:v1.0.5
```

**Staging (همان Production):**
```bash
export ARAS_AUTH_VERSION=v1.1.0
docker-compose up
# استفاده می‌کند از: ghcr.io/aras-group-co/aras-auth:v1.1.0
```

## مزایا

1. **اطلاعات نسخه درون ایمیج**: مقادیر در زمان build در ایمیج قرار می‌گیرند
2. **انعطاف‌پذیری**: امکان تنظیم مقادیر از طریق ARGها
3. **خودکارسازی**: امکان تولید خودکار مقادیر از git
4. **سازگاری**: حفظ سازگاری با روش‌های مختلف build
5. **Multi-tag support**: ساخت همزمان چندین tag (latest, version, local)
6. **Registry ready**: آماده برای push به GitHub Container Registry
7. **Production-safe**: default به version مشخص است
8. **Environment-specific**: هر محیط می‌تواند version خود را تعیین کند
9. **Reproducible**: امکان تکرار دقیق environment در هر زمان
