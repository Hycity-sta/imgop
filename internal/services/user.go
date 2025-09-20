package services

import (
	"imgop/internal/models"
)

func CreateUser(user *models.User) error {
	return models.InsertUser(user)
}
