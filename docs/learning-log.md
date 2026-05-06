# BookHub Learning Log

Bu dokuman, `v0.1-memory-crud` milestone'ina kadar gelinen sureci egitim notu olarak ozetler. Gecmis kod halleri birebir korunmadigi icin eski ara adimlar sahte branch'ler halinde uretilmedi; mevcut calisan CRUD hali tek bir stabil baslangic noktasi olarak kaydedildi.

## 1. Proje Kurulumu

Proje Go module yapisi ile baslatildi. Kodlar, ileride buyumeye uygun olacak sekilde `cmd/server` ve `internal` klasorleri altinda konumlandirildi.

- `cmd/server`: Uygulamanin calistirilabilir HTTP server giris noktasi.
- `internal/models`: Domain model tanimlari.
- `internal/repositories`: Veri erisim soyutlamalari ve memory implementasyonu.

## 2. Ilk HTTP Server

Ilk calisan backend parcasi `net/http` paketi ile kuruldu. Server `http.ListenAndServe(":8080", nil)` ile 8080 portunda baslatildi.

Baslangicta uygulamanin ayakta oldugunu anlamak icin basit handler'lar eklendi:

- `/`
- `/health`
- `/about`

## 3. Handler ve Route Mantigi

HTTP endpoint'leri `http.HandleFunc` ile handler fonksiyonlarina baglandi. Handler fonksiyonlari request'i okuyup response uretmekten sorumlu oldu.

`/books` endpoint'i ayni path uzerinden farkli HTTP method'lari okuyarak CRUD operasyonlarini yonetir:

- `GET`
- `POST`
- `PUT`
- `DELETE`

## 4. Book Modeli

Kitap verisi `Book` modeli ile temsil edildi.

```go
type Book struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
    Year   int    `json:"year"`
}
```

JSON tag'leri sayesinde Go struct alanlari HTTP response ve request body icinde beklenen JSON alan adlariyla eslesti.

## 5. Repository Pattern

Veri erisim mantigini handler'lardan ayirmak icin repository pattern kullanildi. Bu yaklasim, ileride memory repository yerine database repository eklemeyi kolaylastirir.

`BookRepository` interface'i su operasyonlari tanimlar:

- `FindAll`
- `FindByID`
- `Create`
- `Update`
- `Delete`

## 6. Memory Repository

`MemoryBookRepository`, kitaplari uygulama bellegindeki slice icinde tutar. Baslangic verileri repository olusturulurken eklenir ve yeni kitaplar icin `nextID` degeri kullanilir.

Bu yapi kalici veri saklamaz. Server yeniden baslatildiginda eklenen, guncellenen veya silinen veriler sifirlanir.

## 7. Query Parameter Kullanimi

Tek kitap islemleri icin `id` query parameter'i kullanildi.

Ornek:

```text
GET /books?id=1
```

Handler icinde bu deger `r.URL.Query().Get("id")` ile okundu ve `strconv.Atoi` ile `int` tipine cevrildi.

## 8. JSON Response

Response formatinin JSON oldugunu belirtmek icin `Content-Type` header'i ayarlandi:

```go
w.Header().Set("Content-Type", "application/json")
```

Go verileri HTTP response'a `json.NewEncoder(w).Encode(...)` ile yazildi.

## 9. POST ile Create

`POST /books` endpoint'i request body icindeki JSON verisini okuyarak yeni kitap olusturur.

Temel validasyonlar eklendi:

- `title` bos olamaz.
- `author` bos olamaz.

Basarili create isleminde `201 Created` status code'u kullanildi.

## 10. PUT ile Update

`PUT /books?id=1` endpoint'i ID ile mevcut kitabi bulur ve request body'den gelen yeni bilgilerle gunceller.

Guncelleme sirasinda:

- ID query parameter'dan okunur.
- Body JSON olarak decode edilir.
- `title` ve `author` validasyonu yapilir.
- Kayit yoksa `404 Not Found` doner.

## 11. DELETE ile Delete

`DELETE /books?id=1` endpoint'i ID ile eslesen kitabi memory repository icinden siler.

Basarili silme isleminde response body dondurulmez ve `204 No Content` status code'u kullanilir.

## 12. HTTP Status Code Kullanimi

CRUD akisi boyunca temel HTTP status code'lari kullanildi:

- `200 OK`: Basarili listeleme, detay ve guncelleme cevaplari.
- `201 Created`: Basarili kitap olusturma.
- `204 No Content`: Basarili kitap silme.
- `400 Bad Request`: Eksik veya gecersiz input.
- `404 Not Found`: Kitap bulunamadiginda.
- `405 Method Not Allowed`: Desteklenmeyen HTTP method'u kullanildiginda.
- `500 Internal Server Error`: Beklenmeyen repository hatalarinda.

## 13. Interface ve Implementasyon Uyumu

Go'da bir struct'in bir interface'i implement etmesi icin interface'teki tum method'lara sahip olmasi gerekir.

Ornegin `BookRepository` interface'ine `Update` veya `Delete` method'u eklendiginde, `MemoryBookRepository` struct'i da ayni imzaya sahip method'lari implement etmelidir. Aksi halde `MemoryBookRepository`, `BookRepository` olarak kullanilamaz ve compile hatasi olusur.

Bu kural, interface ve implementasyon arasindaki sozlesmeyi guclu tutar.
