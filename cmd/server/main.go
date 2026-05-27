package main

import (
	"BookHub/internal/activity"
	"BookHub/internal/config"
	"BookHub/internal/database"
	"BookHub/internal/middleware"
	"BookHub/internal/mongodb"
	"BookHub/internal/redisdb"
	"BookHub/internal/services"
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"

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
	cfg := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		fmt.Println("Database bağlantısı hatası: ", err)
		return
	}
	defer db.Close()

	err = database.Migrate(ctx, db)
	if err != nil {
		fmt.Println("Migration hatası:", err)
		return
	}

	redisClient, err := redisdb.Connect(ctx, cfg.RedisAddr)
	if err != nil {
		fmt.Println("Redis bağlantısı hatası:", err)
		return
	}
	defer redisClient.Close()

	mongoClient, mongoDB, err := mongodb.Connect(ctx, cfg.MongoURI, cfg.MongoDBName)
	if err != nil {
		fmt.Println("MongoDB bağlantısı hatası:", err)
		return
	}
	defer mongoClient.Disconnect(context.Background())

	mongoActivityLogger := activity.NewMongoActivityLogService(mongoDB)
	activityLogger := activity.NewAsyncActivityLogger(mongoActivityLogger, 100)
	activityLogger.Start()
	defer activityLogger.Stop()

	activityHandler := handlers.NewActivityHandler(activityLogger)

	bookRepo := repositories.NewPostgresBookRepository(db)
	bookService := services.NewBookService(bookRepo, activityLogger)
	bookHandler := handlers.NewBookHandler(bookService)

	userRepo := repositories.NewPostgresUserRepository(db)
	userService := services.NewUserService(userRepo, activityLogger)
	sessionService := services.NewRedisSessionService(redisClient)
	authHandler := handlers.NewAuthHandler(userService, sessionService, activityLogger)
	authMiddleware := middleware.NewAuthMiddleware(sessionService)

	webHandler := handlers.NewWebHandler(bookService, userService, sessionService, activityLogger)

	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/about", aboutHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/books", authMiddleware.RequireAuthForMethods(
		bookHandler.BooksHandler,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
	))
	mux.HandleFunc("/web/books", webHandler.BooksPageHandler)
	mux.HandleFunc("/web/books/new", authMiddleware.RequireWebAuth(webHandler.BookCreatePageHandler))
	mux.HandleFunc("/web/books/edit", authMiddleware.RequireWebAuth(webHandler.BookEditPageHandler))
	mux.HandleFunc("/web/books/delete", authMiddleware.RequireWebAuth(webHandler.BookDeletePageHandler))
	mux.HandleFunc("/web/login", webHandler.LoginPageHandler)
	mux.HandleFunc("/web/register", webHandler.RegisterPageHandler)
	mux.HandleFunc("/web/logout", webHandler.LogoutPageHandler)
	mux.HandleFunc("/register", authHandler.RegisterHandler)
	mux.HandleFunc("/login", authHandler.LoginHandler)
	mux.HandleFunc("/me", authMiddleware.RequireAuth(authHandler.MeHandler))
	mux.HandleFunc("/logout", authMiddleware.RequireAuth(authHandler.LogoutHandler))
	mux.HandleFunc("/activity-logs", authMiddleware.RequireAuth(activityHandler.ListActivityLogsHandler))

	server := &http.Server{
		Addr:         cfg.ServerAddr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		fmt.Println("Server çalışıyor:", serverURL(cfg.ServerAddr))
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Println("Server başlatılamadı:", err)
			stop()
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = server.Shutdown(shutdownCtx)
	if err != nil {
		fmt.Println("Server düzgün kapatılamadı:", err)
	}
}

func serverURL(addr string) string {
	if strings.HasPrefix(addr, ":") {
		return "http://localhost" + addr
	}

	return "http://" + addr
}
