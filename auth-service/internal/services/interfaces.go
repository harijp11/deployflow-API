package services

import "context"

// External api interfaces
type UserCreator interface {
    CreateUser( ctx context.Context,registrationID string,profile UserProfile) (uint, error)
}

type UserProfile struct {
    FirstName   string  `json:"first_name"`
    LastName    string  `json:"last_name"`
    Email       string  `json:"email"`
    CompanyName *string `json:"company_name"`
    PhoneNumber *string `json:"phone_number"`
}