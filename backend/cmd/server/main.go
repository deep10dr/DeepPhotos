package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"deepphotos/backend/internal/config"
	"deepphotos/backend/internal/database"
	"deepphotos/backend/internal/handler"
	appmiddleware "deepphotos/backend/internal/middleware"
	"deepphotos/backend/internal/repository"
	"deepphotos/backend/internal/service"
	"deepphotos/backend/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	log.Println("Initializing DeepPhotos Go Backend Microservice...")

	cfg := config.LoadConfig()

	// 1. Initialize SQLite Database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Fatal: Database initialization failed: %v", err)
	}
	defer db.Close()

	// 2. Initialize MinIO Object Storage Client
	storageClient, err := storage.InitMinio(cfg)
	if err != nil {
		log.Printf("Warning: MinIO storage client init warning: %v", err)
	}

	// 3. Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	photoRepo := repository.NewPhotoRepository(db)
	albumRepo := repository.NewAlbumRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	lockedRepo := repository.NewLockedRepository(db)

	// 4. Initialize Services
	authService := service.NewAuthService(userRepo, auditRepo, cfg.JWTSecret)
	photoService := service.NewPhotoService(photoRepo, storageClient)

	// 5. Initialize Handlers
	authMiddleware := appmiddleware.NewAuthMiddleware(authService)
	authHandler := handler.NewAuthHandler(authService)
	photosHandler := handler.NewPhotosHandler(photoRepo, photoService, storageClient)
	albumsHandler := handler.NewAlbumsHandler(albumRepo)
	usersHandler := handler.NewUsersHandler(userRepo)
	auditHandler := handler.NewAuditHandler(auditRepo)
	lockedHandler := handler.NewLockedHandler(lockedRepo)

	// 6. Router & CORS Setup
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Disable HTTP Caching for API responses to fix UI reactivity issues
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
			next.ServeHTTP(w, r)
		})
	})

	// API Routes Setup
	r.Get("/api/health", handler.HealthHandler)

	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/login", authHandler.Login)
	})

	// Protected Routes
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.RequireAuth)

		r.Route("/api/photos", func(r chi.Router) {
			r.Get("/", photosHandler.List)
			r.Post("/", photosHandler.Upload)
			r.Post("/upload-url", photosHandler.UploadURL)
			r.Post("/batch-delete", photosHandler.DeleteBatch)
			r.Post("/batch-restore", photosHandler.RestoreBatch)
			r.Get("/{id}", photosHandler.GetDetail)
			r.Put("/{id}", photosHandler.Update)
			r.Delete("/{id}", photosHandler.DeleteSingle)
			r.Get("/{id}/file", photosHandler.StreamFile)
			r.Get("/{id}/thumbnail", photosHandler.StreamThumbnail)
		})

		r.Route("/api/albums", func(r chi.Router) {
			r.Get("/", albumsHandler.List)
			r.Post("/", albumsHandler.Create)
			r.Get("/{id}", albumsHandler.GetDetail)
			r.Put("/{id}", albumsHandler.Update)
			r.Delete("/{id}", albumsHandler.Delete)
			r.Post("/{id}/photos", albumsHandler.AddPhotos)
		})

		r.Route("/api/locked-folders", func(r chi.Router) {
			r.Get("/", lockedHandler.List)
			r.Post("/", lockedHandler.Create)
			r.Post("/{id}/verify", lockedHandler.VerifyPasscode)
			r.Delete("/{id}", lockedHandler.Delete)
		})

		r.Route("/api/users", func(r chi.Router) {
			r.Put("/{id}", usersHandler.Update)
			r.Put("/{id}/password", usersHandler.ChangePassword)
			
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireAdmin)
				r.Get("/", usersHandler.List)
				r.Post("/", usersHandler.Create)
				r.Put("/{id}/role", usersHandler.ChangeRole)
				r.Delete("/{id}", usersHandler.Delete)
			})
		})

		r.With(authMiddleware.RequireAdmin).Get("/api/audit-logs", auditHandler.List)
	})

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Graceful shutdown channel
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("DeepPhotos REST API Server listening on http://localhost:%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down REST API server gracefully...")
}
