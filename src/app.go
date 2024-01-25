package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/weldonkipchirchir/go/bookshelf-api/src/config"
	"github.com/weldonkipchirchir/go/bookshelf-api/src/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// port := os.Getenv("PORT")

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	killSig := make(chan os.Signal)

	// MongoDB Connection
	config.DbConnection()

	srv := &http.Server{
		Addr:    ":4000",
		Handler: router,
	}

	routes.SetUpBooks(router)
	routes.SetUpReviews(router)
	routes.SetUpUsers(router)

	router.GET("/healthcheck", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "UP"})
	})

	go func() {
		signal.Notify(killSig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
		<-killSig

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*5)
		defer shutdownCancel()

		go func() {
			<-shutdownCtx.Done()

			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatalln("Server gracefully stopped")
			}
		}()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Fatalln("Error on server shutdown: ", err)
		} else {
			log.Println("Server gracefully stopped")
		}

		serverStopCtx()

	}()

	go func() {
		log.Println("Server is running on http://localhost:4000/")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Close MongoDB connection when the server stops
	defer func() {
		config.CloseDB()
	}()

	<-serverCtx.Done()
}
