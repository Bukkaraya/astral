package activities

// Activities provides payment activity implementations
type Activities struct {
	// In a real implementation, these would be your payment gateway clients,
	// database connections, etc.
}

// NewActivities creates a new payment activities service
func NewActivities() *Activities {
	return &Activities{}
}