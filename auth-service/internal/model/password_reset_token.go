package model

import "time"

type PasswordResetToken struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	CredentialID uint       `gorm:"not null" json:"credential_id"`
	TokenHash    string     `gorm:"type:varchar(64);not null" json:"-"`
	ExpiresAt    time.Time  `gorm:"type:timestamptz;not null" json:"expires_at"`
	UsedAt       *time.Time `gorm:"type:timestamptz" json:"used_at,omitempty"`
	CreatedAt    time.Time  `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}