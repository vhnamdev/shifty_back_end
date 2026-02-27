package main

import (
	"log"
	"os"
	"os/signal"
	"shifty-backend/configs"
	"shifty-backend/graph"
	"shifty-backend/internal/delivery/graphql"
	"shifty-backend/internal/delivery/http/handler"
	"shifty-backend/internal/delivery/http/middleware"
	"shifty-backend/internal/delivery/http/route"
	"shifty-backend/internal/repository"
	"shifty-backend/internal/usecase"
	"shifty-backend/pkg/database"
	"shifty-backend/pkg/mailer"
	"shifty-backend/pkg/monitoring"
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

	if cfg.SentryDSN != "" {
		err := monitoring.Init(cfg.SentryDSN, cfg.AppEnv, cfg.SentryTraceRate)
		if err != nil {
			log.Printf("[SENTRY] Sentry initialization failed: %v", err)
		}
	}
	defer monitoring.Flush()

	// Connect to Redis Database
	redisClient := database.ConnectRedis(cfg)

	// Run Auto Migrate
	database.RunAutoMigrate(db)

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

	app.Use(recover.New())                                 // Auto restart server
	app.Use(logger.New())                                  // Log request to console
	app.Use(middleware.NewRateLimiter(100, 1*time.Minute)) // Rate limit

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
	transactor := repository.NewTransactor(db)
	redisRepo := repository.NewRedisRepo(redisClient)
	userRepo := repository.NewUserRepository(db)
	userRestaurantRepo := repository.NewUserRestaurantRepository(db)
	restaurantRepo := repository.NewRestaurantRepository(db)
	positionRepo := repository.NewPositionRepository(db)
	// ------------------------------USECASE----------------------------------

	authUseCase := usecase.NewAuthUseCase(userRepo, tokenMaster, timeoutContext, redisRepo, emailService, googleService)
	userUseCase := usecase.NewUserUseCase(userRepo, userRestaurantRepo, cloudinaryService, transactor, restaurantRepo)
	userRestaurantUseCase := usecase.NewUserRestaurantUseCase(userRestaurantRepo)
	restaurantUseCase := usecase.NewRestaurantUseCase(transactor, restaurantRepo, userRestaurantRepo, positionRepo, redisRepo, userRepo, emailService, cloudinaryService)
	positionUseCase := usecase.NewPositionUseCase(positionRepo, userRestaurantRepo, transactor)
	// ------------------------------HANDLER----------------------------------

	authHandler := handler.NewAuthHandler(authUseCase, cloudinaryService, emailService)
	userHandler := handler.NewUserHandler(userUseCase, cloudinaryService)
	handlers := &route.AppHandlers{
		AuthHandler: authHandler,
		UserHandler: userHandler,
	}
	gqlResolver := &graph.Resolver{
		UserUseCase:           userUseCase,
		UserRestaurantUseCase: userRestaurantUseCase,
		RestaurantUseCase:     restaurantUseCase,
		PositionUseCase:       positionUseCase,
	}

	playgroundHandler, queryHandler := graphql.NewGraphQLHandler(gqlResolver)
	// Setup routes
	route.SetupRoutes(app, handlers, tokenMaster, playgroundHandler, queryHandler)

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
