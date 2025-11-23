package main

import (
	"github.com/mikhaylovilya/pr-review-service/core/endpoints"
	"github.com/mikhaylovilya/pr-review-service/core/repository"
)

func main() {
	var cache repository.InMemoryRepository = repository.NewInMemory()
	repo := endpoints.NewRepository(&cache)
	server := NewServer(repo)
	server.StartServer()
	// StartServer(&ServerConfig{Addr: ":3081"})
}
