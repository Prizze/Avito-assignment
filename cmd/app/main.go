package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"pr-reviewer/internal/api"
	prDelivery "pr-reviewer/internal/delivery/http/PullRequest"
	teamDelivery "pr-reviewer/internal/delivery/http/Team"
	userDelivery "pr-reviewer/internal/delivery/http/User"
	"pr-reviewer/internal/delivery/http/server"
	"pr-reviewer/internal/pkg/db/postgres"
	"pr-reviewer/internal/pkg/logger"
	"pr-reviewer/internal/pkg/middleware"
	prRepo "pr-reviewer/internal/repository/PullRequest"
	teamRepo "pr-reviewer/internal/repository/Team"
	userRepo "pr-reviewer/internal/repository/User"
	prUC "pr-reviewer/internal/usecase/PullRequest"
	teamUC "pr-reviewer/internal/usecase/Team"
	userUC "pr-reviewer/internal/usecase/User"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../../.env")
}

func main() {
	// Логгер
	l, err := logger.NewZapLogger("warn")
	if err != nil {
		log.Fatalf("failed to make logger: %v", err)
	}

	// Инициализация подключения к БД
	pool, err := postgres.NewPool()
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer pool.Close()

	// Team
	teamRepo := teamRepo.NewTeamRepository(pool)
	teamUC := teamUC.NewTeamUsecase(teamRepo, l)
	teamHandler := teamDelivery.NewTeamHandler(teamUC)

	// User
	userRepo := userRepo.NewUserRepository(pool)
	userUC := userUC.NewUserUsecase(userRepo, l)
	userHandler := userDelivery.NewUserHandler(userUC)

	// PullRequest
	prRepo := prRepo.NewPullRequestRepository(pool)
	prUC := prUC.NewPullRequestUsecase(prRepo, userRepo, l)
	prHandler := prDelivery.NewPRHandler(prUC)

	// Композиция handlers
	server := server.NewServer(userHandler, teamHandler, prHandler)

	r := mux.NewRouter()
	h := api.HandlerWithOptions(server, api.GorillaServerOptions{
		BaseRouter:  r,
		Middlewares: []api.MiddlewareFunc{middleware.RecoverMiddleware},
	})

	addr := ":8080"
	srv := &http.Server{
		Addr:    addr,
		Handler: h,
	}

	// Канал для ловли сигналов остановки
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Println("server started at", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	<-stop
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server stopped gracefully")
}
