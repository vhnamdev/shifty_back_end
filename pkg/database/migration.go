package database

import (
	"log"
	"shifty-backend/internal/entity"

	"gorm.io/gorm"
)

func RunAutoMigrate(db *gorm.DB) {
	log.Println("Loading Auto Migrations")

	err := db.AutoMigrate(
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
		&entity.UserRestaurant{},
	)

	if err != nil {
		log.Fatal("Database Migration Failed: ", err)
	}

	log.Println("Database Migration Successfully!")
}
