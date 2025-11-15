package main

import (
	"log"
	"net/http"
	"pr-reviewer/internal/api"
	pullrequest "pr-reviewer/internal/delivery/http/PullRequest"
	teamDelivery "pr-reviewer/internal/delivery/http/Team"
	user "pr-reviewer/internal/delivery/http/User"
	"pr-reviewer/internal/delivery/http/server"
	"pr-reviewer/internal/pkg/db/postgres"
	teamRepo "pr-reviewer/internal/repository/Team"
	teamUC "pr-reviewer/internal/usecase/Team"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../../.env")
}

func main() {
	pool, err := postgres.NewPool()
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer pool.Close()

	teamRepo := teamRepo.NewTeamRepository(pool)
	teamUC := teamUC.NewTeamUsecase(teamRepo)
	teamHandler := teamDelivery.NewTeamHandler(teamUC)

	userHandler := user.NewUserHandler()
	prHandler := pullrequest.NewPRHandler()

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
