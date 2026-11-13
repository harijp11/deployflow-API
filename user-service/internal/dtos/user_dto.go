package dtos

type CreateUserRequestDTO struct {
    FirstName   string  `json:"first_name" binding:"required,min=2,max=100"`
    LastName    string  `json:"last_name" binding:"required,min=1,max=100"`
    Email       string  `json:"email" binding:"required,email"`
    CompanyName *string `json:"company_name"`
    PhoneNumber *string `json:"phone_number"`
}

type CreateUserResponseDTO struct {
    ID uint `json:"id"`
}