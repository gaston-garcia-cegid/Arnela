package domain

import (
	"time"

	"github.com/google/uuid"
)

// SearchResults represents the aggregated results from global search
type SearchResults struct {
	Clients      []SearchClient      `json:"clients"`
	Employees    []SearchEmployee    `json:"employees"`
	Appointments []SearchAppointment `json:"appointments"`
	Invoices     []SearchInvoice     `json:"invoices"`
	TotalResults int                 `json:"totalResults"`
}

// SearchClient represents a client in search results
type SearchClient struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone,omitempty"`
	DNICIF    string    `json:"dniCif,omitempty"`
}

// SearchEmployee represents an employee in search results
type SearchEmployee struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone,omitempty"`
	Specialties []string  `json:"specialties"`
	AvatarColor string    `json:"avatarColor"`
}

// SearchAppointment represents an appointment in search results
type SearchAppointment struct {
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	StartTime    time.Time `json:"startTime"`
	EndTime      time.Time `json:"endTime"`
	Status       string    `json:"status"`
	ClientName   string    `json:"clientName"`
	EmployeeName string    `json:"employeeName"`
}

// SearchInvoice represents an invoice in search results
type SearchInvoice struct {
	ID            uuid.UUID `json:"id"`
	InvoiceNumber string    `json:"invoiceNumber"`
	ClientName    string    `json:"clientName"`
	TotalAmount   float64   `json:"totalAmount"`
	Status        string    `json:"status"`
	IssueDate     time.Time `json:"issueDate"`
}

// SearchService defines the interface for search operations
type SearchService interface {
	GlobalSearch(query string) (*SearchResults, error)
}
