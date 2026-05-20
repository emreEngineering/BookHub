# API Örnekleri

Bu dosya, BookHub API endpointleri için daha detaylı curl örnekleri içerir. Örnekler uygulamanın `http://localhost:8080` adresinde çalıştığını varsayar.

## Register

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Ada Lovelace","email":"ada@example.com","password":"secret123"}'
```

## Login

Login sonucunda dönen `session_id` cookie değerini `cookies.txt` dosyasına kaydedin:

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{"email":"ada@example.com","password":"secret123"}'
```

## Me

Kaydedilen cookie ile aktif kullanıcıyı görüntüleyin:

```bash
curl http://localhost:8080/me \
  -b cookies.txt
```

## Logout

```bash
curl -X POST http://localhost:8080/logout \
  -b cookies.txt \
  -c cookies.txt
```

## Books CRUD

Kitapları listeleme:

```bash
curl http://localhost:8080/books
```

ID ile kitap getirme:

```bash
curl "http://localhost:8080/books?id=1"
```

Kitap oluşturma:

```bash
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{"title":"Dune","author":"Frank Herbert","year":1965}'
```

Kitap güncelleme:

```bash
curl -X PUT "http://localhost:8080/books?id=1" \
  -H "Content-Type: application/json" \
  -d '{"title":"Dune Messiah","author":"Frank Herbert","year":1969}'
```

Kitap silme:

```bash
curl -X DELETE "http://localhost:8080/books?id=1"
```

## Activity Logs

Activity log kayıtları auth middleware arkasındadır. Önce login olup cookie kaydedin, sonra logları listeleyin:

```bash
curl http://localhost:8080/activity-logs \
  -b cookies.txt
```
