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

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/nickson01/todo/auth"
	"github.com/nickson01/todo/todo"
)

func main() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Printf("Please consider environment variable %s", err)
	}

	var (
		buildCommit = "dev"
		buildTime   = time.Now().String()
	)

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Fail to connect DB")
	}

	db.AutoMigrate(&todo.Todo{})

	r := gin.Default()
	protected := r.Group("", auth.Protect([]byte(os.Getenv("SIGN"))))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/x", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"buildCommit": buildCommit,
			"buildTime":   buildTime,
		})
	})

	r.GET("/tokenz", auth.AccessToken(os.Getenv("SIGN")))

	todoHandler := todo.NewTodoHandler(db)
	protected.POST("/todos", todoHandler.NewtTask)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	s := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s \n", err)
		}
	}()

	<-ctx.Done()
	stop()

	fmt.Println("Shuting down Gracefully, Press Ctrl+c again to force")
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(timeoutCtx); err != nil {
		fmt.Println(err)
	}

}
