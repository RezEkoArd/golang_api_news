package app

import (
	"bwanews/config"

	"bwanews/internal/adapter/cloudflare"
	"bwanews/internal/adapter/handler"
	"bwanews/internal/adapter/repository"
	"bwanews/internal/core/service"
	"bwanews/lib/auth"
	"bwanews/lib/middleware"
	"bwanews/lib/pagination"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"log"
)

func RunServer() {
	cfg := config.NewConfig()
	db, err := cfg.ConnectionPostgres()
	if err != nil {
		log.Fatal("Error connection to database: %v", err)
		return
	}


	//? membuat directory temporary
	err = os.MkdirAll("./temp/content",0755)
	if err != nil {
		log.Fatal("Error creating temp directory: %v", err)
		return
	}


	// Cloudflare
	cdfR2 := cfg.LoadAwsConfig()
	s3Client := s3.NewFromConfig(cdfR2)
	r2Adapter := cloudflare.NewCloudflareR2Adapter(s3Client, cfg)

	jwt := auth.NewJwt(cfg)
	middlewareAuth := middleware.NewMiddleware(cfg) // Middleware

	_ = pagination.NewPagination()

	//? Repository
	authrepo := repository.NewAuthRepository(db.DB)
	categoryRepo := repository.NewCategoryRepository(db.DB)
	contentRepo := repository.NewContentRepository(db.DB)
	userRepo := repository.NewUserRepository(db.DB)
	

	//? Service
	authService := service.NewAuthService(authrepo, cfg, jwt) 
	categoryService := service.NewCategoryService(categoryRepo)
	contentService := service.NewContentService(contentRepo, cfg, r2Adapter)
	userService := service.NewUserService(userRepo)

	//? Handler
	authHandler :=  handler.NewAuthHandler(authService)
	categoryHandler :=  handler.NewCategoryHandler(categoryService)
	contentHandler := handler.NewContentHandler(contentService)
	userHandler := handler.NewUserHandler(userService)

	//? initialize Serve
	app := fiber.New()
	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${ip} ${status} - ${latency} ${method} ${path} \n",
	}))

	// check ENV PROD OR DEV
	if os.Getenv("APP_ENV") != "production" { //"production"|| "development"
		cfg := swagger.Config{
			BasePath: "/",
			FilePath: "./docs/swagger.json",
			Path: "docs",
			Title: "Swagger API Docs",
		}

		app.Use(swagger.New(cfg))
	} 

	api := app.Group("/api")
	api.Post("/login", authHandler.Login)

	// check ENV PROD OR DEV
	if os.Getenv("APP_ENV") != "production" { //"production"|| "development"
		cfg := swagger.Config{
			BasePath: "/api",
			FilePath: "./docs/swagger.json",
			Path: "docs",
			Title: "Swagger API Docs",
		}

		app.Use(swagger.New(cfg))
	} 

	
	adminApp := api.Group("/admin")
	adminApp.Use(middlewareAuth.CheckToken())

	//? Category
	categoryApp := adminApp.Group("/categories")
	categoryApp.Get("/", categoryHandler.GetCategories)
	categoryApp.Post("/", categoryHandler.CreateCategory)
	categoryApp.Put("/:categoryID", categoryHandler.EditCategoryByID)
	categoryApp.Get("/:categoryID", categoryHandler.GetCategoryByID)
	categoryApp.Delete("/:categoryID", categoryHandler.DeleteCategory)


	//? Content
	contentApp := adminApp.Group("/contents")
	contentApp.Get("/",contentHandler.GetContents)
	contentApp.Post("/", contentHandler.CreateContent)
	contentApp.Put("/:contentID", contentHandler.UpdateContent)
	contentApp.Get("/:contentID", contentHandler.GetContentByID)
	contentApp.Delete("/:contentID", contentHandler.DeleteContent)
	contentApp.Post("/upload-image", contentHandler.UploadImageR2)

	//? User
	userApp := adminApp.Group("/users/profile")
	userApp.Get("/",userHandler.GetUserByID)
	userApp.Put("/update-password",userHandler.UpdatePassword)
 

	//? FE
	feApp := api.Group("/fe")
	feApp.Get("/categories", categoryHandler.GetCategoryFE)
	feApp.Get("/contents",contentHandler.GetContentWithQuery)
	feApp.Get("/contents/:contentID", contentHandler.GetContentDetail)

	//? routes
	 go func() {
		if cfg.App.AppPort == "" {
			cfg.App.AppPort = os.Getenv("APP_PORT")
		}

		err := app.Listen(":" + cfg.App.AppPort)
		if err!= nil {
			log.Fatal("Error starting server: %v", err)
		} 
	 }()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	<-quit

	log.Println("server shutdown of 5 seconds")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	app.ShutdownWithContext(ctx)
	
}