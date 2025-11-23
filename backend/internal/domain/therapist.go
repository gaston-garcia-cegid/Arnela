package domain

// Therapist represents a therapist/psychologist in the system
// This is a mock implementation until the employees feature is complete
type Therapist struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Specialties []string `json:"specialties"`
	IsAvailable bool     `json:"isAvailable"`
	AvatarColor string   `json:"avatarColor"` // For UI display
}

// GetMockTherapists returns the three mock therapists
func GetMockTherapists() []Therapist {
	return []Therapist{
		{
			ID:          "therapist-1",
			Name:        "Dra. Alicia Mendoza",
			Specialties: []string{"Terapia Infantil", "Terapia Familiar"},
			IsAvailable: true,
			AvatarColor: "#6366F1", // primary color
		},
		{
			ID:          "therapist-2",
			Name:        "Dra. Sandra Martínez",
			Specialties: []string{"Terapia para Adolescentes", "Terapia para Adultos"},
			IsAvailable: true,
			AvatarColor: "#10B981", // secondary color
		},
		{
			ID:          "therapist-3",
			Name:        "Dra. Guadalupe Torres",
			Specialties: []string{"Terapia para Adolescentes", "Terapia para Adultos"},
			IsAvailable: true,
			AvatarColor: "#F59E0B", // accent color
		},
	}
}

// GetTherapistByID finds a therapist by ID from mock data
func GetTherapistByID(id string) *Therapist {
	therapists := GetMockTherapists()
	for _, t := range therapists {
		if t.ID == id {
			return &t
		}
	}
	return nil
}

// IsValidTherapistID checks if a therapist ID is valid
func IsValidTherapistID(id string) bool {
	return GetTherapistByID(id) != nil
}
