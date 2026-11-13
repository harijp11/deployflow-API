package models

import "time"

type RegistrationIntent struct {
	ID                    uint       `gorm:"primaryKey;autoIncrement"`
	FirstName             string     `gorm:"type:varchar(100);not null"`
	LastName              string     `gorm:"type:varchar(100);not null"`
	Email                 string     `gorm:"type:varchar(255);not null;uniqueIndex"`
	PasswordHash          string     `gorm:"type:varchar(255);not null"`
	CompanyName           *string    `gorm:"type:varchar(255)"`
	PhoneNumber           *string    `gorm:"type:varchar(20)"`
	VerificationTokenHash string     `gorm:"type:varchar(64);not null;uniqueIndex"`
	ExpiresAt             time.Time  `gorm:"type:timestamp;not null;index"`
	EmailSendStatus       string     `gorm:"type:varchar(20);not null;default:pending"`
	EmailSendAttempts     int        `gorm:"not null;default:0"`
	LastEmailSentAt       *time.Time `gorm:"type:timestamp"`
	CreatedAt             time.Time  `gorm:"type:timestamp;not null;autoCreateTime"`
	UpdatedAt             time.Time  `gorm:"type:timestamp;not null;autoUpdateTime"`
}

func (RegistrationIntent) TableName() string {
	return "registration_intents"
}