package initializers

import "github.com/yonraz/gochat_auth/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}