package dtos

type CreateUserRequestDTO struct {
    FirstName   string  `json:"first_name" binding:"required,min=2,max=100"`
    LastName    string  `json:"last_name" binding:"required,min=1,max=100"`
    Email       string  `json:"email" binding:"required,email"`
    CompanyName *string `json:"company_name"`
    PhoneNumber *string `json:"phone_number"`
}
type UpdateUserRequestDTO struct {
    FirstName   *string `json:"first_name" binding:"omitempty,min=2,max=100"`
    LastName    *string `json:"last_name" binding:"omitempty,min=1,max=100"`
    CompanyName *string `json:"company_name" binding:"omitempty,max=150"`
    PhoneNumber *string `json:"phone_number" binding:"omitempty,max=20"`
}

type UserResponseDTO struct {
    ID          uint    `json:"id"`
    FirstName   string  `json:"first_name"`
    LastName    string  `json:"last_name"`
    Email       string  `json:"email"`
    CompanyName *string `json:"company_name,omitempty"`
    PhoneNumber *string `json:"phone_number,omitempty"`
}

type CreateUserResponseDTO struct {
    ID uint `json:"id"`
}