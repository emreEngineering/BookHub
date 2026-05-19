package main

import (
	"BookHub/internal/database"
	"BookHub/internal/gormdb"
	"BookHub/internal/middleware"
	"BookHub/internal/services"
	"fmt"
	"net/http"

	"BookHub/internal/handlers"
	"BookHub/internal/repositories"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "BookHub çalışıyor")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "BookHub bir kitap yönetim sistemidir")
}

func main() {
	db, err := database.Connect()
	if err != nil {
		fmt.Println("Database bağlantısı hatası: ", err)
		return
	}
	defer db.Close()

	err = database.Migrate(db)
	if err != nil {
		fmt.Println("Migration hatası:", err)
		return
	}

	gormDB, err := gormdb.Connect()
	if err != nil {
		fmt.Println("GORM database bağlantısı hatası:", err)
		return
	}

	err = gormdb.AutoMigrate(gormDB)
	if err != nil {
		fmt.Println("GORM migration hatası:", err)
		return
	}

	bookRepo := repositories.NewGormBookRepository(gormDB)
	bookService := services.NewBookService(bookRepo)
	bookHandler := handlers.NewBookHandler(bookService)

	userRepo := repositories.NewGormUserRepository(gormDB)
	userService := services.NewUserService(userRepo)
	sessionService := services.NewSessionService()
	authHandler := handlers.NewAuthHandler(userService, sessionService)
	authMiddleware := middleware.NewAuthMiddleware(sessionService)

	webHandler := handlers.NewWebHandler(bookService, userService, sessionService)

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/about", aboutHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/books", bookHandler.BooksHandler)
	http.HandleFunc("/web/books", webHandler.BooksPageHandler)
	http.HandleFunc("/web/books/new", webHandler.BookCreatePageHandler)
	http.HandleFunc("/web/books/edit", webHandler.BookEditPageHandler)
	http.HandleFunc("/web/books/delete", webHandler.BookDeletePageHandler)
	http.HandleFunc("/web/login", webHandler.LoginPageHandler)
	http.HandleFunc("/web/register", webHandler.RegisterPageHandler)
	http.HandleFunc("/web/logout", webHandler.LogoutPageHandler)
	http.HandleFunc("/register", authHandler.RegisterHandler)
	http.HandleFunc("/login", authHandler.LoginHandler)
	http.HandleFunc("/me", authMiddleware.RequireAuth(authHandler.MeHandler))
	http.HandleFunc("/logout", authMiddleware.RequireAuth(authHandler.LogoutHandler))
	fmt.Println("Server çalışıyor: http://localhost:8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server başlatılamadı:", err)
	}
}
