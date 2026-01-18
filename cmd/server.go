package main

import (
	"log"
	"os"
	"os/signal"
	"shifty-backend/configs"
	handler "shifty-backend/internal/delivery/http"
	"shifty-backend/internal/delivery/http/route"
	"shifty-backend/internal/domain"
	"shifty-backend/pkg/database"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load config file
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatal("Can not load config: ", err)
	}
	log.Println("Config loaded successfully!")

	// Connect to PostgreSQL Database
	db := database.ConnectPostgres(cfg)
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Can not connect to PostgreSQL Database!")
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Fatal("Can not disconnect PostgresSQL Database!")
		}
	}()
	// Run Auto Migrate
	log.Println("Loading Auto Migrations")
	err = db.AutoMigrate(
		&domain.Restaurant{},
		&domain.Position{},
		&domain.User{},
		&domain.ShiftRule{},
		&domain.Schedule{},
		&domain.Shift{},
		&domain.ShiftRequirement{},
		&domain.ShiftRequest{},
		&domain.ShiftAssignment{},
		&domain.Post{},
		&domain.Comment{},
		&domain.Reaction{},
		&domain.Conversation{},
		&domain.Participant{},
		&domain.Feedback{},
	)

	if err != nil {
		log.Fatal("Database Migration Failed: ", err)
	}
	log.Println("Database Migration Successfully!")

	app := fiber.New(fiber.Config{
		AppName:      "Shifty Backend API",
		ErrorHandler: handler.GlobalErrorHandler,
	})

	app.Use(recover.New()) // Auto restart server
	app.Use(logger.New())  // Log request to console
	// Config Cors
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, HEAD, PUT, DELETE, PATCH",
	}))
	handlers := &route.AppHandlers{}
	route.SetupRoutes(app, handlers)
	go func() {
		port := cfg.AppPort
		if port == "" {
			port = ":8080"
		} else {
			if port[0] != ':' {
				port = ":" + port
			}
		}
		log.Printf(" Server starting on port %s\n", port)
		if err := app.Listen(port); err != nil {
			log.Panic(err)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	log.Println("Shutting down server")

	_ = app.Shutdown()
	log.Println("Server exited successfully!")
}
