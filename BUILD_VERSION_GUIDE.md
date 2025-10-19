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

## Single Clean Version Strategy

### فلسفه نهایی

**فقط ایمیج‌های با Git Tag تمیز قابل Push هستند**

- ✅ فقط version تمیز (مثلاً `v1.3.0`)
- ❌ بدون `latest` در هیچ‌جا
- ❌ بدون `-dirty` در push
- ❌ بدون `-N-gXXXXXX` در push
- ✅ Dev و Production دقیقاً همان ایمیج

### چرا این استراتژی؟

❌ **مشکلات `latest` و dirty versions:**
- **غیرقابل پیش‌بینی**: نمی‌دانید دقیقاً کدام version را اجرا می‌کنید
- **عدم امکان Rollback**: اگر مشکلی پیش آید، نمی‌توانید به version قبلی برگردید
- **ناسازگاری بین سرورها**: سرورهای مختلف ممکن است version های متفاوت داشته باشند
- **مشکلات Debugging**: نمی‌توانید دقیقاً بگویید کدام version مشکل دارد
- **Breaking Changes**: اگر version جدیدی push شود، ممکن است سیستم شما خراب شود
- **Overwrite خطرناک**: dirty versions ممکن است clean versions را overwrite کنند

### استراتژی پیاده‌سازی شده

✅ **همه محیط‌ها (docker-compose.yml و docker-compose.dev.yml):**
```yaml
services:
  aras_auth:
    image: ghcr.io/aras-group-co/aras-auth:${ARAS_AUTH_VERSION}
```
- **بدون default value**: باید صریحاً version را set کنید
- **یکسان در همه جا**: Dev و Production دقیقاً همان version

### نحوه استفاده

**Development (local testing):**
```bash
# Build هر version ای (حتی dirty)
make docker-build-versioned
# نتیجه: ghcr.io/aras-group-co/aras-auth:v1.2.0-dirty

# تست محلی
export ARAS_AUTH_VERSION=v1.2.0-dirty
docker compose up
```

**Release به Production:**
```bash
# 1. اطمینان از clean state
git status

# 2. Commit همه چیز
git add .
git commit -m "feat: ready for v1.3.0"

# 3. ایجاد tag
git tag v1.3.0
git push origin v1.3.0

# 4. Build
make docker-build-versioned
# نتیجه: ghcr.io/aras-group-co/aras-auth:v1.3.0

# 5. Push (فقط الان کار می‌کند!)
make docker-push-versioned
# ✅ Success: Pushing v1.3.0

# 6. Deploy در همه محیط‌ها
export ARAS_AUTH_VERSION=v1.3.0
docker compose up  # Production
docker compose -f docker-compose.dev.yml up  # Development
```

### Validation و Error Messages

اگر سعی کنید dirty version را push کنید:

```bash
make docker-push-versioned

# ❌ ERROR: Cannot push non-release version: v1.2.0-dirty
# 
# Only clean git tags can be pushed to registry.
# 
# Current issues:
#   - You have uncommitted changes
# 
# To fix:
#   1. Commit all changes: git add . && git commit
#   2. Create a tag: git tag v1.x.x
#   3. Push tag: git push origin v1.x.x
```

## مزایا

1. **اطلاعات نسخه درون ایمیج**: مقادیر در زمان build در ایمیج قرار می‌گیرند
2. **امنیت کامل**: فقط release های رسمی در registry
3. **سادگی**: یک tag، یک ایمیج
4. **Validation**: خطای واضح با راهنمای رفع
5. **Consistency**: Dev و Prod همیشه یکسان
6. **No surprises**: نمی‌توان اشتباهی dirty push کرد
7. **Traceability**: هر ایمیج به یک git tag مشخص اشاره دارد
8. **Reproducible**: امکان تکرار دقیق environment در هر زمان
9. **Clean builds**: فقط از committed code ایمیج production ساخته می‌شود
10. **Registry ready**: آماده برای push به GitHub Container Registry
