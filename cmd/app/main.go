package main

import (
	"log"
	"net/http"
	"pr-reviewer/internal/api"
	prDelivery "pr-reviewer/internal/delivery/http/PullRequest"
	teamDelivery "pr-reviewer/internal/delivery/http/Team"
	userDelivery "pr-reviewer/internal/delivery/http/User"
	"pr-reviewer/internal/delivery/http/server"
	"pr-reviewer/internal/pkg/db/postgres"
	"pr-reviewer/internal/pkg/logger"
	prRepo "pr-reviewer/internal/repository/PullRequest"
	teamRepo "pr-reviewer/internal/repository/Team"
	userRepo "pr-reviewer/internal/repository/User"
	prUC "pr-reviewer/internal/usecase/PullRequest"
	teamUC "pr-reviewer/internal/usecase/Team"
	userUC "pr-reviewer/internal/usecase/User"

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

	teamRepo := teamRepo.NewTeamRepository(pool)
	teamUC := teamUC.NewTeamUsecase(teamRepo, l)
	teamHandler := teamDelivery.NewTeamHandler(teamUC)

	userRepo := userRepo.NewUserRepository(pool)
	userUC := userUC.NewUserUsecase(userRepo, l)
	userHandler := userDelivery.NewUserHandler(userUC)

	prRepo := prRepo.NewPullRequestRepository(pool)
	prUC := prUC.NewPullRequestUsecase(prRepo, userRepo, l)
	prHandler := prDelivery.NewPRHandler(prUC)

	server := server.NewServer(userHandler, teamHandler, prHandler)

	r := mux.NewRouter()
	h := api.HandlerWithOptions(server, api.GorillaServerOptions{
		BaseRouter: r,
	})

	addr := ":8080"
	log.Println("server started")
	if err := http.ListenAndServe(addr, h); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
