package generator

import (
	"clientgen-v2/internal/models"
)

// Generator defines the interface for code generation
type Generator interface {
	Generate(data *models.TemplateData) ([]byte, error)
	SetTemplate(template Template) error
	SetFormatter(formatter Formatter) error
}

// Template defines template operations
type Template interface {
	Render(data interface{}) (string, error)
	Load(templatePath string) error
	AddFunction(name string, fn interface{}) error
}

// Formatter handles code formatting
type Formatter interface {
	Format(code []byte) ([]byte, error)
	Validate(code []byte) error
}

// GenerationStrategy defines different generation approaches
type GenerationStrategy interface {
	ShouldGenerate(method *models.WorkflowMethod) bool
	GetTemplateName() string
	TransformData(data *models.TemplateData) *models.TemplateData
}