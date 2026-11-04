
package model

import (
	"time"

	"github.com/google/uuid"
)

type AuthSession struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	CredentialID uuid.UUID  `gorm:"type:uuid;not null" json:"credential_id"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	TokenHash    string     `gorm:"type:varchar(255);not null;unique" json:"-"`
	FamilyID     uuid.UUID  `gorm:"type:uuid;not null" json:"family_id"`
	ExpiresAt    time.Time  `gorm:"type:timestamptz;not null" json:"expires_at"`
	RevokedAt    *time.Time `gorm:"type:timestamptz" json:"revoked_at,omitempty"`
	LastUsedAt   *time.Time `gorm:"type:timestamptz" json:"last_used_at,omitempty"`
	CreatedAt    time.Time  `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"type:timestamptz;not null" json:"updated_at"`
}

func (AuthSession) TableName() string {
	return "auth_sessions"
}
