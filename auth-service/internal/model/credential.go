
package model

import (
	"time"

	"github.com/google/uuid"
)

type Credential struct {
	ID                uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	UserID            uuid.UUID  `gorm:"type:uuid;not null;unique" json:"user_id"`
	Email             string     `gorm:"type:varchar(320);not null" json:"email"`
	PasswordHash      string     `gorm:"type:text;not null" json:"-"`
	EmailVerifiedAt   *time.Time `gorm:"type:timestamptz" json:"email_verified_at,omitempty"`
	PasswordChangedAt time.Time  `gorm:"type:timestamptz;not null" json:"password_changed_at"`
	CreatedAt         time.Time  `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"type:timestamptz;not null" json:"updated_at"`
}

func (Credential) TableName() string {
	return "credentials"
}
