package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/Nishal77/resona/backend/internal/auth"
	"github.com/Nishal77/resona/backend/internal/comments"
	"github.com/Nishal77/resona/backend/internal/communities"
	"github.com/Nishal77/resona/backend/internal/engagements"
	"github.com/Nishal77/resona/backend/internal/explore"
	"github.com/Nishal77/resona/backend/internal/language"
	"github.com/Nishal77/resona/backend/internal/notifications"
	"github.com/Nishal77/resona/backend/internal/posts"
	"github.com/Nishal77/resona/backend/internal/upload"
	"github.com/Nishal77/resona/backend/internal/users"
	"github.com/Nishal77/resona/backend/internal/vrs"
	"github.com/Nishal77/resona/backend/pkg/config"
	"github.com/Nishal77/resona/backend/pkg/database"
	rdb "github.com/Nishal77/resona/backend/pkg/redis"
	"github.com/Nishal77/resona/backend/pkg/supabase"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Config
	config.Load()

	// Infrastructure
	database.Connect(config.App.DatabaseURL)
	database.Migrate()
	rdb.Connect(config.App.RedisURL)

	// Supabase Storage
	storage := supabase.NewStorageClient(
		config.App.SupabaseURL,
		config.App.SupabaseServiceKey,
		config.App.SupabaseStorageBucket,
	)

	// Services (dependency order)
	langSvc := language.NewService()
	vrsSvc := vrs.NewService()
	engSvc := engagements.NewService()
	notifSvc := notifications.NewService()

	usersRepo := users.NewRepository()
	postRepo := posts.NewRepository()

	postSvc := posts.NewService(postRepo, langSvc, vrsSvc, engSvc, notifSvc, usersRepo)
	usersSvc := users.NewService(usersRepo)

	// Handlers
	authHandler := auth.NewHandler(auth.NewService(auth.NewRepository()))
	usersHandler := users.NewHandler(usersSvc)
	postsHandler := posts.NewHandler(postSvc)
	commentsHandler := comments.NewHandler(comments.NewRepository(), notifSvc)
	communitiesHandler := communities.NewHandler(communities.NewRepository())
	notifHandler := notifications.NewHandler(notifSvc)
	exploreHandler := explore.NewHandler()
	uploadHandler := upload.NewHandler(storage)

	// VRS Scheduler
	scheduler := vrs.NewScheduler(vrsSvc)
	scheduler.Start()
	defer scheduler.Stop()

	// Gin
	if config.App.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{config.App.FrontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true, // required for httpOnly cookie
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	authHandler.Register(api)
	usersHandler.Register(api)
	postsHandler.Register(api)
	commentsHandler.Register(api)
	communitiesHandler.Register(api)
	notifHandler.Register(api)
	exploreHandler.Register(api)
	uploadHandler.Register(api)

	log.Info().Str("port", config.App.Port).Msg("resona backend starting")

	go func() {
		if err := r.Run(":" + config.App.Port); err != nil {
			log.Fatal().Err(err).Msg("server failed")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down")
}
