package domain

// DashboardStats represents the statistics displayed on the dashboard
type DashboardStats struct {
	Clients      ClientStats      `json:"clients"`
	Employees    EmployeeStats    `json:"employees"`
	Appointments AppointmentStats `json:"appointments"`
}

// ClientStats represents client-related statistics
type ClientStats struct {
	Total    int `json:"total"`
	Active   int `json:"active"`
	Inactive int `json:"inactive"`
}

// EmployeeStats represents employee-related statistics
type EmployeeStats struct {
	Total    int `json:"total"`
	Active   int `json:"active"`
	Inactive int `json:"inactive"`
}

// AppointmentStats represents appointment-related statistics
type AppointmentStats struct {
	Total     int `json:"total"`
	Pending   int `json:"pending"`
	Confirmed int `json:"confirmed"`
	Completed int `json:"completed"`
	Cancelled int `json:"cancelled"`
}
