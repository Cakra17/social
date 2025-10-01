package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cakra17/social/internal/handlers"
	"github.com/cakra17/social/internal/store"
	"github.com/cakra17/social/pkg/jwt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Use(middleware.Timeout(time.Minute))

	server := http.Server{
		Addr: ":6969",
		Handler: r,
		WriteTimeout: 30 * time.Second,
		ReadTimeout: 10 * time.Second,
		IdleTimeout: time.Minute,
	}
	// connection
	db := store.ConnectDB(store.DBConfig{
		DB_USERNAME: "admin",
		DB_PASSWORD: "adminsecret",
		DB_HOST: "localhost",
		DB_PORT: "5432",
		DB_NAME: "social",
		DB_MaxOpenConn: 30,
		DB_MaxIdleConn: 30,
		DB_MaxConnLifetime: 15 * time.Minute,
		DB_MaxConnIdletime: 15 * time.Minute,
	})

	jwtAuthenticator := jwt.NewJWTAuthenticator("mysecret")

	// repository
	userRepo := store.NewUserRepo(db)
	postRepo := store.NewPostRepo(db)

	// handler
	userHandler := handlers.NewUserHandler(handlers.UserHandlerConfig{
		UserRepo: userRepo,
		JWTAuthenticator: jwtAuthenticator,
	})

	posthandler := handlers.NewPostHandler(handlers.PostHandlerConfig{
		PostRepo: postRepo,
	})

	posthandler.Init()
	// routing
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/login", userHandler.Authenticate)
		// user
		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.CreateUser)

			r.Group(func(r chi.Router) {
				r.Use(jwtAuthenticator.JWTMiddleware)
				r.Put("/{id}", userHandler.UpdateUser)
				r.Delete("/{id}", userHandler.DeleteUser)
			})
		})
		

		// post
		r.Route("/posts", func(r chi.Router) {
			r.Use(jwtAuthenticator.JWTMiddleware)
			r.Post("/", posthandler.CreatePost)
			r.Put("/{id}", posthandler.UpdatePost)
			r.Delete("/{id}", posthandler.DeletePost)
		})
	})

	ctx := context.Background()
	closed := make(chan struct{})

	// gracefully shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		signal := <-sigint

		log.Printf("Received %s signal, shutting down server", signal.String())
		ctx, cancel := context.WithTimeout(ctx, 5 * time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}
		close(closed)
	}()

	log.Printf("server running on port %s", server.Addr[1:])
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Failed to run server: %v", err)
	}

	<-closed
	log.Println("Server shutdown gracefully")
}