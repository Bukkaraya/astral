package generator

import (
	"fmt"
	"go/format"
	"text/template"
	"strings"
	"time"

	"clientgen-v2/internal/config"
	"clientgen-v2/internal/models"
)

// ClientGenerator generates client code for Temporal workflows
type ClientGenerator struct {
	template  Template
	formatter Formatter
	strategy  GenerationStrategy
	config    *config.GeneratorConfig
}

// NewClientGenerator creates a new client generator
func NewClientGenerator(config *config.GeneratorConfig) *ClientGenerator {
	return &ClientGenerator{
		template:  NewGoTemplate(),
		formatter: NewGoFormatter(),
		strategy:  NewClientStrategy(),
		config:    config,
	}
}

// Generate generates client code from template data
func (g *ClientGenerator) Generate(data *models.TemplateData) ([]byte, error) {
	// Enrich template data with generation metadata
	if data.Metadata == nil {
		data.Metadata = &models.GenerationMetadata{}
	}
	data.Metadata.GeneratedAt = time.Now().Format(time.RFC3339)
	data.Metadata.Generator = "clientgen-v2"
	data.Metadata.Version = "2.0.0"

	// Apply strategy transformation
	if g.strategy != nil {
		data = g.strategy.TransformData(data)
	}

	// Render template
	rendered, err := g.template.Render(data)
	if err != nil {
		return nil, fmt.Errorf("failed to render template: %w", err)
	}

	// Format code
	formatted, err := g.formatter.Format([]byte(rendered))
	if err != nil {
		return nil, fmt.Errorf("failed to format generated code: %w", err)
	}

	return formatted, nil
}

// SetTemplate sets the template
func (g *ClientGenerator) SetTemplate(template Template) error {
	g.template = template
	return nil
}

// SetFormatter sets the formatter
func (g *ClientGenerator) SetFormatter(formatter Formatter) error {
	g.formatter = formatter
	return nil
}

// GoTemplate implements Template interface using Go's text/template
type GoTemplate struct {
	template *template.Template
	funcMap  template.FuncMap
}

// NewGoTemplate creates a new Go template
func NewGoTemplate() *GoTemplate {
	return &GoTemplate{
		funcMap: getTemplateFunctions(),
	}
}

// Render renders the template with data
func (t *GoTemplate) Render(data interface{}) (string, error) {
	if t.template == nil {
		return "", fmt.Errorf("template not loaded")
	}

	var buf strings.Builder
	if err := t.template.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// Load loads template from string  
func (t *GoTemplate) Load(templatePath string) error {
	tmpl, err := template.New("client").Funcs(t.funcMap).Parse(templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}
	t.template = tmpl
	return nil
}

// AddFunction adds a custom function to the template
func (t *GoTemplate) AddFunction(name string, fn interface{}) error {
	t.funcMap[name] = fn
	return nil
}

// GoFormatter implements Formatter interface using go/format
type GoFormatter struct{}

// NewGoFormatter creates a new Go formatter
func NewGoFormatter() *GoFormatter {
	return &GoFormatter{}
}

// Format formats Go code
func (f *GoFormatter) Format(code []byte) ([]byte, error) {
	formatted, err := format.Source(code)
	if err != nil {
		return nil, fmt.Errorf("failed to format Go code: %w", err)
	}
	return formatted, nil
}

// Validate validates Go code syntax
func (f *GoFormatter) Validate(code []byte) error {
	_, err := format.Source(code)
	return err
}

// ClientStrategy implements GenerationStrategy for client generation
type ClientStrategy struct{}

// NewClientStrategy creates a new client strategy
func NewClientStrategy() *ClientStrategy {
	return &ClientStrategy{}
}

// ShouldGenerate determines if a method should generate client code
func (s *ClientStrategy) ShouldGenerate(method *models.WorkflowMethod) bool {
	return method != nil && method.Name != ""
}

// GetTemplateName returns the template name for client generation
func (s *ClientStrategy) GetTemplateName() string {
	return "client"
}

// TransformData transforms data for client generation
func (s *ClientStrategy) TransformData(data *models.TemplateData) *models.TemplateData {
	// Add client-specific transformations if needed
	return data
}

// getTemplateFunctions returns template functions
func getTemplateFunctions() template.FuncMap {
	return template.FuncMap{
		"toJSONTag":     toJSONTag,
		"toKebabCase":   toKebabCase,
		"extractIDField": extractIDField,
		"extractEntity": extractEntity,
		"toDescription": toDescription,
		"toLower":       strings.ToLower,
		"ne":           func(a, b string) bool { return a != b },
	}
}

// Template helper functions
func toJSONTag(field string) string {
	if len(field) == 0 {
		return field
	}
	return strings.ToLower(string(field[0])) + field[1:]
}

func toKebabCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('-')
		}
		result.WriteRune(rune(strings.ToLower(string(r))[0]))
	}
	return result.String()
}

func extractIDField(inputType string) string {
	entityName := strings.TrimSuffix(inputType, "Request")
	if entityName != "" && entityName != inputType {
		return entityName + "ID"
	}
	return "ID"
}

func extractEntity(inputType string) string {
	entityName := strings.TrimSuffix(inputType, "Request")
	if entityName != "" && entityName != inputType {
		return strings.ToLower(entityName)
	}
	return "entity"
}

func toDescription(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune(' ')
		}
		if i == 0 {
			result.WriteRune(r)
		} else {
			result.WriteRune(rune(strings.ToLower(string(r))[0]))
		}
	}
	return result.String()
}