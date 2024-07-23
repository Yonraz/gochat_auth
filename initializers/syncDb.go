package initializers

import "github.com/yonraz/gochat_auth/models"

func SyncDatabase() {
	db.AutoMigrate(&models.User{})
}