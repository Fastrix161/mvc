package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fastrix161/mvc/pkg/api"
	"github.com/fastrix161/mvc/pkg/models"
)

func main() {
	_, err := models.InitDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	router := api.SetupRouter()
	
	server := &http.Server{
		Addr:    ":8100",
		Handler: router,
	}
	go func() {
		fmt.Printf("Starting server on %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}
	if err := models.CloseDatabase(); err != nil {
		log.Printf("Error closing database: %v", err)
	}

	fmt.Println("Server exited gracefully")
}
