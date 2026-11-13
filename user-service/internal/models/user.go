package models

import "time"

type User struct {
    ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    RegistrationID string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"-"`
    FirstName      string    `gorm:"type:varchar(100);not null" json:"first_name"`
    LastName       string    `gorm:"type:varchar(100);not null" json:"last_name"`
    Email          string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
    CompanyName    *string   `gorm:"type:varchar(255)" json:"company_name"`
    PhoneNumber    *string   `gorm:"type:varchar(20)" json:"phone_number"`
    CreatedAt      time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
    UpdatedAt      time.Time `gorm:"not null;autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
    return "users"
}