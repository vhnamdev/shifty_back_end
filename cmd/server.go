package main

import (
	"log"
	"os"
	"os/signal"
	"shifty-backend/configs"
	"shifty-backend/internal/delivery/http/handler"
	"shifty-backend/internal/delivery/http/route"
	"shifty-backend/internal/entity"
	"shifty-backend/internal/repository"
	"shifty-backend/internal/usecase"
	"shifty-backend/pkg/database"
	"shifty-backend/pkg/mailer"
	"shifty-backend/pkg/token"
	"shifty-backend/pkg/uploader"
	"syscall"
	"time"

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

	// Connect to Redis Database
	redisClient := database.ConnectRedis(cfg)
	// Run Auto Migrate
	log.Println("Loading Auto Migrations")
	err = db.AutoMigrate(
		&entity.Restaurant{},
		&entity.Position{},
		&entity.User{},
		&entity.ShiftRule{},
		&entity.Schedule{},
		&entity.Shift{},
		&entity.ShiftRequirement{},
		&entity.ShiftRequest{},
		&entity.ShiftAssignment{},
		&entity.Post{},
		&entity.Comment{},
		&entity.Reaction{},
		&entity.Conversation{},
		&entity.Participant{},
		&entity.Feedback{},
	)
	if err != nil {
		log.Fatal("Database Migration Failed: ", err)
	}
	log.Println("Database Migration Successfully!")

	accessDuration, err := time.ParseDuration(cfg.AcessTokenExpiry)
	if err != nil {
		log.Fatal("Invalid access token duration")
	}

	refreshDuration, err := time.ParseDuration(cfg.RefreshTokenExpiry)
	if err != nil {
		log.Fatal("Invalid refresh token duration")
	}
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

	// Token
	tokenMaster := token.NewToken(
		cfg.JWTAccessSecret,
		cfg.JWTRefreshSecret,
		accessDuration,
		refreshDuration,
	)

	// Set default timeout if missing
	timeoutInt := cfg.ContextTimeout
	if timeoutInt <= 0 {
		timeoutInt = 5 // Default 5 seconds
	}
	timeoutContext := time.Duration(timeoutInt) * time.Second
	// ----------------------- INFRASTRUCTURE--------------------------------

	emailService := mailer.NewGoMail(cfg.SMTPHost, cfg.SMTPPort, cfg.GmailUser, cfg.GmailPassword)
	cloudinaryService, err := uploader.NewCloudinary(cfg.CloudName, cfg.CloudinaryAPIKey, cfg.CloudinaryAPISecret, cfg.CloudinaryFolderName)
	googleService := configs.NewGoogleConfig(cfg.GoogleClientID, cfg.GoogleSecret, "postmessage")

	// -----------------------REPOSITORY-------------------------------------
	redisRepo := repository.NewRedisRepo(redisClient)
	userRepo := repository.NewUserRepository(db)

	// ------------------------------USECASE----------------------------------

	authUseCase := usecase.NewAuthUseCase(userRepo, tokenMaster, timeoutContext, redisRepo, emailService, googleService)

	// ------------------------------HANDLER----------------------------------

	authHandler := handler.NewAuthHandler(authUseCase, cloudinaryService, emailService)

	handlers := &route.AppHandlers{
		AuthHandler: authHandler,
	}

	// Setup routes
	route.SetupRoutes(app, handlers, tokenMaster)

	// Start server
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
