package repository

import (
	"context"

	"gorm.io/gorm"
)

type ChangePasswordRepository interface {
	UpdatePassword(ctx context.Context, db *gorm.DB, username string, hashedPassword string) error
}

type changePasswordRepository struct{}

func NewChangePasswordRepository() *changePasswordRepository {
	return &changePasswordRepository{}
}

func (repo *changePasswordRepository) UpdatePassword(ctx context.Context, db *gorm.DB, username string, hashedPassword string) error {
	return db.WithContext(ctx).Table("users").Where("username = ?", username).Update("password", hashedPassword).Error
}
