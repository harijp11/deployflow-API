package services

import "context"

// External api interfaces
type UserCreator interface {
    CreateUser( ctx context.Context,registrationID string,profile UserProfile) (uint, error)
}

type UserProfile struct {
    FirstName   string
    LastName    string
    Email       string
    CompanyName *string
    PhoneNumber *string
}