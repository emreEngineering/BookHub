package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Connect fonksiyonu oluşturuluyor
// Fonksiyon parametreleri:
// 1. *sql.DB -> veritabanı bağlantısı
// 2. error   -> hata bilgisi
func Connect() (*sql.DB, error) {

	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env dosyası yüklenemedi, ortam değişkenleri kullanılacak")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL bulunamadı")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("PostgreSQL bağlantısı başarılı")
	return db, nil
}

/*
 Connect fonksiyonunun genel çalışma mantığı:

 Bu fonksiyon uygulamanın PostgreSQL veritabanına bağlanmasını sağlar.
 İlk olarak .env dosyasını okumaya çalışır.
 Daha sonra DATABASE_URL değişkenini alır.

 Eğer DATABASE_URL bulunamazsa hata döndürür.
 Çünkü veritabanının adresi olmadan bağlantı kurulamaz.

 Ardından PostgreSQL bağlantısı açılır.
 Açılan bağlantının gerçekten çalışıp çalışmadığını kontrol etmek için
 db.Ping() kullanılır.

 Eğer herhangi bir aşamada hata oluşursa:
 return nil, err
 ile fonksiyon durur ve hata geri gönderilir.

 Eğer tüm işlemler başarılıysa:
 çalışan db bağlantısı geri döndürülür.

*/
