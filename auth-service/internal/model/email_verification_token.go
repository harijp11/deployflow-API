
package model

import "time"

type EmailVerificationToken struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	CredentialID uint       `gorm:"not null;index" json:"credential_id"`
	TokenHash    string     `gorm:"type:varchar(64);not null;unique" json:"-"`
	ExpiresAt    time.Time  `gorm:"type:timestamptz;not null" json:"expires_at"`
	UsedAt       *time.Time `gorm:"type:timestamptz" json:"used_at"`
	CreatedAt    time.Time  `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (EmailVerificationToken) TableName() string {
	return "email_verification_tokens"
}
