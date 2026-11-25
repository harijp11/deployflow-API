package utils

// ResourceNotFoundError represents a missing resource.
type ResourceNotFoundError struct {
	Message string
}

func (e *ResourceNotFoundError) Error() string {
	return e.Message
}

func NewResourceNotFoundError(message string) *ResourceNotFoundError {
	return &ResourceNotFoundError{Message: message}
}

// ValidationError represents invalid field values or business validation.
type ValidationError struct {
	Message string
	Code    string
	Field   string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func NewValidationError(message, code, field string) *ValidationError {
	return &ValidationError{
		Message: message,
		Code:    code,
		Field:   field,
	}
}

// ConflictError represents duplicate or conflicting data.
type ConflictError struct {
	Message string
}

func (e *ConflictError) Error() string {
	return e.Message
}

func NewConflictError(message string) *ConflictError {
	return &ConflictError{Message: message}
}

// UnauthorizedError represents missing or invalid authentication.
type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

func NewUnauthorizedError(message string) *UnauthorizedError {
	return &UnauthorizedError{Message: message}
}

// ForbiddenError represents insufficient permissions.
type ForbiddenError struct {
	Message string
}

func (e *ForbiddenError) Error() string {
	return e.Message
}

func NewForbiddenError(message string) *ForbiddenError {
	return &ForbiddenError{Message: message}
}

// BadRequestError represents a malformed request.
type BadRequestError struct {
	Message string
}

func (e *BadRequestError) Error() string {
	return e.Message
}

func NewBadRequestError(message string) *BadRequestError {
	return &BadRequestError{Message: message}
}

// InternalServerError represents an unexpected server-side failure.
type InternalServerError struct {
	Message string
}

func (e *InternalServerError) Error() string {
	return e.Message
}

func NewInternalServerError(message string) *InternalServerError {
	return &InternalServerError{Message: message}
}