package models

import (
	"time"

	"github.com/google/uuid"
)

type AuthSession struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CredentialID uint      `gorm:"not null" json:"credential_id"`
	UserID       uint      `gorm:"not null" json:"user_id"`
	TokenHash    string    `gorm:"type:varchar(64);not null;unique" json:"-"`
	FamilyID     uuid.UUID `gorm:"type:uuid;not null" json:"-"`

	// Device Information
	DeviceName      *string `gorm:"type:varchar(100)" json:"device_name"`
	DeviceType      *string `gorm:"type:varchar(50)" json:"device_type"`
	Browser         *string `gorm:"type:varchar(100)" json:"browser"`
	OperatingSystem *string `gorm:"type:varchar(100)" json:"operating_system"`
	IPAddress       *string `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent       *string `gorm:"type:text" json:"-"`

	ExpiresAt  time.Time  `gorm:"type:timestamptz;not null" json:"expires_at"`
	RevokedAt  *time.Time `gorm:"type:timestamptz" json:"revoked_at"`
	LastUsedAt *time.Time `gorm:"type:timestamptz" json:"last_used_at"`
	CreatedAt  time.Time  `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"type:timestamptz;not null" json:"updated_at"`
}

func (AuthSession) TableName() string {
	return "auth_sessions"
}
