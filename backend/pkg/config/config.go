package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	AppEnv      string
	FrontendURL string

	DatabaseURL string

	SupabaseURL        string
	SupabaseServiceKey string
	SupabaseStorageBucket string

	RedisURL string

	JWTSecret            string
	JWTRefreshSecret     string
	JWTExpiresIn         time.Duration
	JWTRefreshExpiresIn  time.Duration

	GoogleClientID     string
	GoogleClientSecret string
	GoogleCallbackURL  string

	DeepLAPIKey string

	VRSTrendingThreshold  float64
	VRSShareVelocityHours int
}

var App *Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, reading from environment")
	}

	App = &Config{
		Port:        getEnv("PORT", "8080"),
		AppEnv:      getEnv("APP_ENV", "development"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),

		DatabaseURL: mustGetEnv("DATABASE_URL"),

		SupabaseURL:           mustGetEnv("SUPABASE_URL"),
		SupabaseServiceKey:    mustGetEnv("SUPABASE_SERVICE_KEY"),
		SupabaseStorageBucket: getEnv("SUPABASE_STORAGE_BUCKET", "resona-media"),

		RedisURL: mustGetEnv("REDIS_URL"),

		JWTSecret:           mustGetEnv("JWT_SECRET"),
		JWTRefreshSecret:    mustGetEnv("JWT_REFRESH_SECRET"),
		JWTExpiresIn:        parseDuration("JWT_EXPIRES_IN", 15*time.Minute),
		JWTRefreshExpiresIn: parseDuration("JWT_REFRESH_EXPIRES_IN", 7*24*time.Hour),

		GoogleClientID:     mustGetEnv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: mustGetEnv("GOOGLE_CLIENT_SECRET"),
		GoogleCallbackURL:  mustGetEnv("GOOGLE_CALLBACK_URL"),

		DeepLAPIKey: getEnv("DEEPL_API_KEY", ""),

		VRSTrendingThreshold:  parseFloat("VRS_TRENDING_THRESHOLD", 0.75),
		VRSShareVelocityHours: parseInt("VRS_SHARE_VELOCITY_HOURS", 2),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required env var %s is not set", key)
	}
	return v
}

func parseDuration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}

func parseFloat(key string, fallback float64) float64 {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fallback
	}
	return f
}

func parseInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return i
}
