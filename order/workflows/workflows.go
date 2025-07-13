package workflows

// Workflows provides order workflow implementations
type Workflows struct {
	activities Activities
}

// NewWorkflows creates a new order workflows service
func NewWorkflows(activities Activities) *Workflows {
	return &Workflows{
		activities: activities,
	}
}