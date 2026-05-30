# YAPAY ZEKA İLE ADIM ADIM GO ÖĞRENİYORUM


# BookHub

BookHub, Go ile yazılmış bir kitap yönetim sistemidir. Projede hem JSON API hem de `html/template` ile hazırlanmış Web UI bulunur.

Uygulama; `database/sql` ile açık SQL kullanan PostgreSQL kalıcılığı, Redis tabanlı session storage, MongoDB activity log kayıtları, goroutine/channel ile çalışan async activity worker ve unit testlerle adım adım geliştirilmiş bir öğrenme projesidir.

## Final Release

- Current stable release: `v1.0.0`
- Release notes: [docs/release-notes-v1.0.0.md](docs/release-notes-v1.0.0.md)

## Özellikler

- Book CRUD API
- Web üzerinden kitap listeleme, oluşturma, düzenleme ve silme
- Kullanıcı register/login/logout akışı
- Cookie tabanlı session yönetimi
- Redis session storage
- PostgreSQL persistence with `database/sql`
- MongoDB activity logs
- Goroutine/channel ile async activity worker
- Standart JSON response yapısı
- Service, middleware ve activity worker unit testleri
- Git learning tag ve branch sistemi

## Teknolojiler

- Go
- `net/http`
- `html/template`
- PostgreSQL
- `database/sql`
- Redis
- MongoDB
- bcrypt
- Go testing package
- Git

## Gereksinimler

- Go
- PostgreSQL
- Redis
- MongoDB
- Git

## Kurulum

Repoyu clone edin:

```bash
git clone <repo-url>
cd BookHub
```

Ortam değişkenleri için örnek dosyayı kopyalayın:

```bash
cp .env.example .env
```

Gerekli environment değişkenleri:

```env
SERVER_ADDR=:8080
DATABASE_URL=postgres://username:password@localhost:5432/bookhub?sslmode=disable
REDIS_ADDR=localhost:6379
MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=bookhub
```

`DATABASE_URL` değerindeki kullanıcı adı, şifre ve veritabanı adını kendi yerel PostgreSQL kurulumunuza göre düzenleyin.

## Servisleri Çalıştırma

PostgreSQL çalışıyor olmalı ve `DATABASE_URL` içinde belirtilen veritabanı erişilebilir durumda olmalıdır.

Redis'i başlatın ve bağlantıyı kontrol edin:

```bash
brew services start redis
redis-cli ping
```

Beklenen Redis cevabı:

```text
PONG
```

MongoDB'yi başlatın ve bağlantıyı kontrol edin:

```bash
brew services start mongodb-community@8.0
mongosh --eval "db.runCommand({ ping: 1 })"
```

## Uygulamayı Çalıştırma

Bağımlılıkları düzenleyin:

```bash
go mod tidy
```

Server'ı başlatın:

```bash
go run ./cmd/server
```

Beklenen başlangıç çıktıları:

- PostgreSQL bağlantısı başarılı
- Redis bağlantısı başarılı
- MongoDB bağlantısı başarılı
- Server çalışıyor

Uygulama varsayılan olarak `http://localhost:8080` adresinde çalışır.

## API Endpointleri

| Method | Endpoint | Açıklama |
| --- | --- | --- |
| `GET` | `/` | Ana sayfa/metin endpointi |
| `GET` | `/health` | Sağlık kontrolü |
| `GET` | `/about` | Proje hakkında kısa bilgi |
| `GET` | `/books` | Kitapları listeler |
| `GET` | `/books?id=1` | Tek kitap getirir |
| `POST` | `/books` | Yeni kitap oluşturur; login gerektirir |
| `PUT` | `/books?id=1` | Kitap günceller; login gerektirir |
| `DELETE` | `/books?id=1` | Kitap siler; login gerektirir |
| `POST` | `/register` | Kullanıcı oluşturur |
| `POST` | `/login` | Giriş yapar ve session cookie döner |
| `GET` | `/me` | Aktif kullanıcı bilgisini getirir |
| `POST` | `/logout` | Session'ı siler |
| `GET` | `/activity-logs` | Activity log kayıtlarını listeler |

## Web Endpointleri

| Method | Endpoint | Açıklama |
| --- | --- | --- |
| `GET` | `/web/books` | Kitap listesi sayfası |
| `GET` | `/web/books/new` | Yeni kitap formu; login gerektirir |
| `POST` | `/web/books/new` | Yeni kitap oluşturma; login gerektirir |
| `GET` | `/web/books/edit?id=1` | Kitap düzenleme formu; login gerektirir |
| `POST` | `/web/books/edit?id=1` | Kitap düzenleme submit; login gerektirir |
| `POST` | `/web/books/delete?id=1` | Kitap silme; login gerektirir |
| `GET` | `/web/login` | Giriş formu |
| `POST` | `/web/login` | Giriş submit |
| `GET` | `/web/register` | Kayıt formu |
| `POST` | `/web/register` | Kayıt submit |
| `POST` | `/web/logout` | Web çıkış işlemi |

## Örnek curl Komutları

Register:

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Ada Lovelace","email":"ada@example.com","password":"secret123"}'
```

Login ve cookie kaydetme:

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{"email":"ada@example.com","password":"secret123"}'
```

Aktif kullanıcı:

```bash
curl http://localhost:8080/me \
  -b cookies.txt
```

Kitap oluşturma:

```bash
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{"title":"Dune","author":"Frank Herbert","year":1965}'
```

Activity logs:

```bash
curl http://localhost:8080/activity-logs \
  -b cookies.txt
```

Daha detaylı örnekler için [docs/api-examples.md](docs/api-examples.md) dosyasına bakabilirsiniz.

## Öğrenme Sırası

Bu projeyi en verimli şekilde öğrenmek için önce basit HTTP akışını, sonra katmanları, en son dış servisleri okuyun. Her adımda şu soruları sorun: Bu dosyanın tek sorumluluğu ne? Bu katman hangi katmana bağımlı? Bu iş burada mı olmalı, yoksa bir alt veya üst katmanda mı?

1. `/books` API akışını inceleyin.
   Başlangıç dosyaları: `cmd/server/main.go`, `internal/handlers/book_handler.go`, `internal/responses/response.go`.

2. Domain modellerini okuyun.
   Başlangıç dosyaları: `internal/models/book.go`, `internal/models/user.go`.

3. Service layer mantığını öğrenin.
   Başlangıç dosyaları: `internal/services/book_services.go`, `internal/services/user_service.go`.

4. Repository ve SQL katmanına geçin.
   Başlangıç dosyaları: `internal/repositories/postgres_book_repository.go`, `internal/repositories/postgres_user_repository.go`, `internal/database/database.go`, `internal/database/migrations.go`.

5. `context.Context` akışını takip edin.
   Akış: `r.Context()` -> service -> repository -> PostgreSQL/Redis/MongoDB.

6. Auth ve session mekanizmasını öğrenin.
   Başlangıç dosyaları: `internal/handlers/auth_handler.go`, `internal/services/redis_session_service.go`, `internal/middleware/auth_middleware.go`, `internal/requestcontext/user.go`.

7. Web template tarafını inceleyin.
   Başlangıç dosyaları: `internal/handlers/web_handler.go`, `templates/`.

8. MongoDB activity log yapısını okuyun.
   Başlangıç dosyaları: `internal/activity/`, `internal/mongodb/mongodb.go`.

9. Async worker, goroutine ve channel kullanımını öğrenin.
   Başlangıç dosyası: `internal/activity/async_activity_logger.go`.

10. Config ve server startup akışını inceleyin.
    Başlangıç dosyaları: `internal/config/config.go`, `cmd/server/main.go`.

11. Testleri katman katman okuyun.
    Başlangıç dosyaları: `internal/services/*_test.go`, `internal/middleware/*_test.go`, `internal/activity/*_test.go`.

## Testler

Tüm testleri çalıştırma:

```bash
go test ./...
```

Servis testleri:

```bash
go test ./internal/services -cover
```

Middleware testleri:

```bash
go test ./internal/middleware -cover
```

Activity worker testleri:

```bash
go test ./internal/activity -cover
```

## Git Öğrenme Sistemi

BookHub, özellikleri küçük adımlarla öğrenmek için tag ve branch düzeniyle tutulur.

- `feature/*` branch'leri yeni özellik geliştirme adımları için kullanılır.
- `learning/*` branch'leri konu anlatımı ve adım adım öğrenme akışları için kullanılır.
- `v0.1-memory-crud-learning-final`, Memory CRUD öğrenme bölümünün audit edilmiş final tag'idir.
- Öğrenme audit notları için [docs/learning-audit.md](docs/learning-audit.md) dosyasına bakılabilir.

## Sürüm Geçmişi

- `v0.1` Memory CRUD
- `v0.2` Handler Layer
- `v0.3` Service Layer
- `v0.4` Standard Response
- `v0.5` Auth
- `v0.6` Template
- `v0.7` Database
- `v0.8` GORM
- `v0.9` Redis
- `v0.10` MongoDB
- `v0.11` Background Worker
- `v0.12` Tests
- `v0.13-gorm-final` Final GORM mimarisi
- `v0.13-sql-final` Final `database/sql` mimarisi

## Güvenlik Notları

- `.env` commit edilmez.
- Gerçek parola, token veya secret değerleri repoya eklenmez.
- Kullanıcı parolaları bcrypt ile hashlenir.
- Session bilgileri Redis'te TTL ile tutulur.
- `cookies.txt` commit edilmez.
