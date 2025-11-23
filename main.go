package main

import (
	"github.com/mikhaylovilya/pr-review-service/core/endpoints"
	"github.com/mikhaylovilya/pr-review-service/core/storage"
)

func main() {
	var cache storage.InMemoryRepository = storage.NewInMemory()
	repo := endpoints.NewRepository(&cache)
	server := NewServer(repo)
	server.StartServer()
	// StartServer(&ServerConfig{Addr: ":3081"})
}
