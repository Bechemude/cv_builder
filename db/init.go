package db

import (
	"cvbuilder/config"
	"cvbuilder/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	Postgres *gorm.DB
}

func Init(c *config.Config) (*DB, error) {
	dsn := "host=" + c.PostgresHost + " port=" + c.PostgresPort + " user=" + c.PostgresUser + " password=" + c.PostgresPass + " dbname=" + c.PostgresDBName + " sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return nil, err
	}

	log.Println("Connected to Postgres database")

	log.Println("Starting database migrations...")
	err = db.AutoMigrate(
		&models.User{},
		&models.CV{},
		&models.JobHistory{},
		&models.Job{},
	)
	if err != nil {
		log.Printf("Migration error: %v", err)
		return nil, err
	}
	log.Println("Database migrations completed successfully")

	return &DB{
		Postgres: db,
	}, nil
}
