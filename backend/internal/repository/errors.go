package repository

import "errors"

// Repository errors
var (
	// Client errors
	ErrClientNotFound = errors.New("client not found")

	// Employee errors
	ErrEmployeeNotFound = errors.New("employee not found")

	// Common errors
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrDNIAlreadyExists   = errors.New("DNI already exists")
	ErrNIFAlreadyExists   = errors.New("NIF already exists")

	// Appointment errors
	ErrAppointmentNotFound = errors.New("appointment not found")

	// User errors
	ErrUserNotFound = errors.New("user not found")
)
