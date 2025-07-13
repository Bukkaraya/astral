package workflows

// Workflows provides payment workflow implementations
type Workflows struct {
	activities Activities
}

// NewWorkflows creates a new payment workflows service
func NewWorkflows(activities Activities) *Workflows {
	return &Workflows{
		activities: activities,
	}
}