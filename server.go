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

	"github.com/gin-gonic/gin"
	"github.com/mikhaylovilya/pr-review-service/core/endpoints"
)

type Server struct {
	Repository *endpoints.Repository
}

func NewServer(repo *endpoints.Repository) Server {
	return Server{
		Repository: repo,
	}
}

func (s *Server) StartServer() {
	router := gin.Default()

	router.POST("/team/add", s.Repository.AddTeamHandler)
	router.GET("/team/get/:teamName", s.Repository.GetTeamHandler)
	// router.POST("/users/setIsActive", s.Repository.AddTeamHandler)
	router.POST("/pullRequest/create", s.Repository.CreatePullRequestHandler)
	router.POST("/pullRequest/reassign", s.Repository.ReassignHandler)

	httpServer := &http.Server{
		Addr:    ":3081",
		Handler: router,
	}

	go func() {
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server fatal error: %v", err)
		}
		log.Println("Stopped recieving new requests...")
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	timeout := 5 * time.Second
	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}
	log.Println("Server shutdown was completed gracefully")
}
