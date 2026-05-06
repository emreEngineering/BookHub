# BookHub

BookHub, Go ile yazılmış küçük bir backend projesidir. Projenin amacı kitap kayıtlarını yönetmek için sade bir HTTP CRUD API geliştirmek ve ilerleyen bölümlerde katmanlı mimari, standart response yapısı, authentication, template, database, cache, concurrency ve test konularını adım adım uygulamaktır.

## Mevcut Ozellikler

- Bellek uzerinde calisan kitap deposu
- `Book` modeli
- `BookRepository` interface'i
- `MemoryBookRepository` implementasyonu
- JSON request ve response destegi
- Query parameter ile kitap ID okuma
- Kitap listeleme, detay goruntuleme, olusturma, guncelleme ve silme
- Temel HTTP status code kullanimi

## Kullanilan Teknolojiler

- Go
- `net/http`
- `encoding/json`
- Git

## Endpointler

| Method | Endpoint | Aciklama |
| --- | --- | --- |
| `GET` | `/books` | Tum kitaplari listeler |
| `GET` | `/books?id=1` | ID degerine gore tek kitap getirir |
| `POST` | `/books` | Yeni kitap olusturur |
| `PUT` | `/books?id=1` | ID degerine gore kitap gunceller |
| `DELETE` | `/books?id=1` | ID degerine gore kitap siler |

## Calistirma

```bash
go run ./cmd/server
```

Server varsayilan olarak `http://localhost:8080` adresinde calisir.

## Test

```bash
go test ./...
```

## curl Ornekleri

Tum kitaplari listeleme:

```bash
curl http://localhost:8080/books
```

ID ile kitap getirme:

```bash
curl "http://localhost:8080/books?id=1"
```

Yeni kitap olusturma:

```bash
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{"title":"Dune","author":"Frank Herbert","year":1965}'
```

Kitap guncelleme:

```bash
curl -X PUT "http://localhost:8080/books?id=1" \
  -H "Content-Type: application/json" \
  -d '{"title":"Suç ve Ceza","author":"Fyodor Dostoyevski","year":1866}'
```

Kitap silme:

```bash
curl -X DELETE "http://localhost:8080/books?id=1"
```

## Git Akisi

Mevcut memory CRUD API hali `v0.1-memory-crud` tag'i ile sabitlenmistir. `main` branch'i stabil milestone'lari tutar. Yeni bolumler icin `feature/*` branch'leri kullanilir ve calisan bolumler tamamlandikca `main` uzerine alinip tag'lenir.

Planlanan branch sirasi:

- `feature/02-handler-layer`
- `feature/03-service-layer`
- `feature/04-standard-response`
- `feature/05-auth`
- `feature/06-template`
- `feature/07-database`
- `feature/08-nosql-cache`
- `feature/09-concurrency`
- `feature/10-tests`
