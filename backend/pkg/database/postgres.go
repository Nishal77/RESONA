package database

import (
	"log"

	"github.com/Nishal77/resona/backend/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(dsn string) {
	var err error
	// PreferSimpleProtocol disables prepared statements.
	// Required for Supabase transaction pooler (port 6543) which doesn't support them.
	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)

	log.Println("database connected")
}

func Migrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Community{},
		&models.CommunityMember{},
		&models.Post{},
		&models.Tag{},
		&models.Comment{},
		&models.Engagement{},
		&models.Follow{},
		&models.Notification{},
		&models.RefreshToken{},
	)
	if err != nil {
		log.Fatalf("auto-migrate failed: %v", err)
	}
	log.Println("database migrated")
}
