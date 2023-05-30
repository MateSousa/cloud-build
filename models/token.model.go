package models

import (
	"time"

	"github.com/MateSousa/cloud-build/initializers"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Token struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID    string    `gorm:"not null" json:"user_id"`
	User      User
	Token     string    `gorm:"not null" json:"token"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"` // soft delete
}

func (t *Token) GenerateToken(user User) (string, error) {

	// load the secret from the config
	config, err := initializers.LoadConfig(".")
	if err != nil {
		return "", err
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign the token with the secret
	tokenString, err := token.SignedString([]byte(config.SecretToken))
	if err != nil {
		return "", err
	}

	// Save the token
	t.UserID = user.ID.String()
	t.Token = tokenString
	t.ExpiresAt = time.Now().Add(time.Hour * 24)
	result := initializers.DB.Create(&t)
	if result.Error != nil {
		return "", result.Error
	}

	return tokenString, nil
}

func (t *Token) ValidateToken(tokenString string) (bool, error) {
	// load the secret from the config
	config, err := initializers.LoadConfig(".")
	if err != nil {
		return false, err
	}

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.SecretToken), nil
	})
	if err != nil {
		return false, err
	}

	// Check if the token is valid
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true, nil
	}

	return false, nil
}

func (t *Token) FindByToken(tokenString string) error {
	result := initializers.DB.Where("token = ?", tokenString).First(&t)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (t *Token) Destroy(tokenString string) error {
	result := t.FindByToken(tokenString)
	if result != nil {
		return result
	}
	err := initializers.DB.Where("token = ?", tokenString).Delete(&t)
	if err.Error != nil {
		return err.Error
	}

	return nil

}
